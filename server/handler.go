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

func (s *Server) handleBlock(block *types.Block) error {
	if err := s.processContract(block.Transactions()); err != nil {
		return err
	}
	return nil
}

func (s *Server) processContract(txs types.Transactions) error {
	for _, tx := range txs {
		ca, err := s.deployTransaction(tx)
		if err != nil {
			continue
		}
		fmt.Println("dectect CA")

		var (
			journal = make([][]byte, 0)
			except  = make(map[interface{}]struct{})
			cas     = make([]common.Address, 1)
		)

		// tx.Data also contains initialization code that will never
		// be used again, we use CodeAt to store the bytecode.
		//
		// tx.Data = initial code + byte code
		cas[0] = ca

		if s.cfg.AllowProxyContract {
			// == Transparent Proxy
			if admin, impl, err := s.eip1967(tx); err == nil {
				cas = append(cas, admin, impl)
			}

			// == UUPS Proxy
			if logic, err := s.eip1822(tx); err == nil {
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
				TxHash:     tx.Hash().Hex(),
				Candidates: r,
			})

			if err != nil {
				// Proxy pattern allows different contracts to point to the
				// same implementation contract, so we ignores 'ErrAlreadyExist'.
				if errors.Is(err, database.ErrAlreadyExist) {
					// This data is saved from a previous request and should not
					// be deleted even if the request fails.
					except[addr.Bytes()] = struct{}{}
					continue
				}

				// Revert; occur database error.
				for _, task := range journal {
					// If a database error occurs, the delete request will likely
					// fail as well, but it doesn't seem to be critical.
					if _, excepted := except[task]; !excepted {
						s.engine.Delete(task)
					}
				}
				return fmt.Errorf("request failed in database: %v", err)
			}

			journal = append(journal, addr.Bytes())
		}
	}

	return nil
}
