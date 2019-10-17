package remouse

//go:generate mockgen -destination mock_driver_test.go -package remouse -self_package github.com/kevinconway/remouse/pkg github.com/kevinconway/remouse/pkg Driver
//go:generate mockgen -destination mock_positionscaler_test.go -package remouse -self_package github.com/kevinconway/remouse/pkg github.com/kevinconway/remouse/pkg PositionScaler
//go:generate mockgen -destination mock_statemachine_test.go -package remouse -self_package github.com/kevinconway/remouse/pkg github.com/kevinconway/remouse/pkg StateMachine
//go:generate mockgen -destination mock_evdeviterator_test.go -package remouse -self_package github.com/kevinconway/remouse/pkg github.com/kevinconway/remouse/pkg EvdevIterator
