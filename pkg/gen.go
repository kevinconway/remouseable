package remouseable

//go:generate mockgen -destination mock_driver_test.go -package remouseable -self_package github.com/kevinconway/remouseable/pkg github.com/kevinconway/remouseable/pkg Driver
//go:generate mockgen -destination mock_positionscaler_test.go -package remouseable -self_package github.com/kevinconway/remouseable/pkg github.com/kevinconway/remouseable/pkg PositionScaler
//go:generate mockgen -destination mock_statemachine_test.go -package remouseable -self_package github.com/kevinconway/remouseable/pkg github.com/kevinconway/remouseable/pkg StateMachine
//go:generate mockgen -destination mock_evdeviterator_test.go -package remouseable -self_package github.com/kevinconway/remouseable/pkg github.com/kevinconway/remouseable/pkg EvdevIterator
//go:generate mockgen -destination mock_readcloser_test.go -package remouseable -self_package github.com/kevinconway/remouseable/pkg io ReadCloser
