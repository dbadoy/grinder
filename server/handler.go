package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/dbadoy/grinder/pkg/database"
	"github.com/dbadoy/grinder/pkg/grinder"
	"github.com/dbadoy/grinder/server/dto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Server) handleBlock(block *types.Block) (err error) {
	defer func() {
		if err != nil {
			s.revert()
		}
	}()

	return s.handleTransactions(block.Transactions())
}

func (s *Server) handleTransactions(txs types.Transactions) (err error) {
	for _, tx := range txs {
		if ca, err := contractAddress(tx); err == nil {
			// Do handleContract if it is a deployment transaction.
			if err := s.handleContract(tx.Hash(), ca); err != nil {
				return err
			}
		}

		/*
			Additional logics
		*/
	}

	return nil
}

func (s *Server) handleContract(hash common.Hash, ca common.Address) error {
	var (
		cas = make([]common.Address, 1)
	)

	// tx.Data also contains initialization code that will never
	// be used again, we use CodeAt to store the bytecode.
	//
	// tx.Data = initial code + byte code
	cas[0] = ca

	if s.cfg.AllowProxyContract {
		// == Transparent Proxy
		if admin, impl, err := s.eip1967(ca); err == nil {
			cas = append(cas, admin, impl)
		}

		// == UUPS Proxy
		if logic, err := s.eip1822(ca); err == nil {
			cas = append(cas, logic)
		}
	}

	for _, addr := range cas {
		code, err := s.eth.CodeAt(context.Background(), addr, nil)
		if err != nil {
			return err
		}

		methods, events, err := grinder.Grinde(code)
		if err != nil {
			return err
		}

		r := make([]string, len(methods)+len(events))
		copy(r[0:], methods)
		copy(r[len(methods):], events)

		err = s.engine.Insert([]byte(addr.Hex()), &dto.Contract{
			TxHash:     hash.Hex(),
			Candidates: r,
		})

		if err != nil {
			// Proxy pattern allows different contracts to point to the
			// same implementation contract, so we ignores 'ErrAlreadyExist'.
			//
			// This data is saved from a previous request and should not be
			// reverted if the request fails. Skip without appending to the
			// 'journals'.
			if errors.Is(err, database.ErrAlreadyExist) {
				continue
			}

			return fmt.Errorf("request failed in database: %v", err)
		}

		s.journals = append(s.journals, &insertContract{[]byte(addr.Hex())})
	}

	return nil
}

func (s *Server) handleRequest(req Request) {
	var err error

	switch req.Kind() {
	case abiRequestType:
		abi := req.(*ABIRequest)
		err = s.engine.Insert([]byte(abi.Name), abi.ABI)

	case contractRequestType:
		contract := req.(*ContractRequest)
		err = s.handleContract(common.Hash{} /* TODO */, contract.Address)

	default:
		err = errors.New("invalid request")
	}

	if err != nil {
		s.revert()
	}

	req.Errorc() <- err
}

// revert performs a revert to a previous state if an
// intermediate failure occurs when making multiple
// requests to the engine within a single request.
func (s *Server) revert() {
	for _, task := range s.journals {
		// If a database error occurs, the delete request will likely
		// fail as well.
		//
		// TODO(dbadoy): We can leave it as a file and perform the
		// Revert when the server is restarted after the DB is
		// recovered.
		task.revert(s.engine)
	}

	if len(s.journals) != 0 {
		s.journals = make([]journalObject, 0)
	}
}
