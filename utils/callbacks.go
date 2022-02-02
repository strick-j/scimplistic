package utils

import (
	"errors"
	"log"
)

const (
	// ErrorIncorrectFunction is thrown when function input or return parameters don't match with the callback
	ErrorIncorrectFunction = "Incorrect function parameters or return parameters"
	// ErrorServerRunning is thrown when an action cannot be taken because the server is running. Pausing the server
	// will enable you to run the command.
	ErrorServerRunning = "Cannot call when the server is running."
)

// SetStartCallback sets the callback that triggers when the server first starts up. The
// function passed must have the same parameter types as the following example:
//
//    func serverStarted(){
//	     //code...
//	 }
func SetStartCallback(cb interface{}) error {
	if serverStarted {
		return errors.New(ErrorServerRunning)
	} else if callback, ok := cb.(func()); ok {
		startCallback = callback
	} else {
		return errors.New(ErrorIncorrectFunction)
	}
	return nil
}

// SetPauseCallback sets the callback that triggers when the server is paused. The
// function passed must have the same parameter types as the following example:
//
//    func serverPaused(){
//	     //code...
//	 }
func SetPauseCallback(cb interface{}) error {
	if serverStarted {
		return errors.New(ErrorServerRunning)
	} else if callback, ok := cb.(func()); ok {
		pauseCallback = callback
	} else {
		return errors.New(ErrorIncorrectFunction)
	}
	return nil
}

// SetResumeCallback sets the callback that triggers when the server is resumed after being paused. The
// function passed must have the same parameter types as the following example:
//
//    func serverResumed(){
//	     //code...
//	 }
func SetResumeCallback(cb interface{}) error {
	if serverStarted {
		return errors.New(ErrorServerRunning)
	} else if callback, ok := cb.(func()); ok {
		resumeCallback = callback
	} else {
		return errors.New(ErrorIncorrectFunction)
	}
	return nil
}

// SetShutDownCallback sets the callback that triggers when the server is shut down. The
// function passed must have the same parameter types as the following example:
//
//    func serverStopped(){
//	     //code...
//	 }
func SetShutDownCallback(cb interface{}) error {
	if serverStarted {
		log.Println("this is the interface 1", cb)
		return errors.New(ErrorServerRunning)
	} else if callback, ok := cb.(func()); ok {
		log.Println("this is the interface", cb)
		stopCallback = callback
	} else {
		return errors.New(ErrorIncorrectFunction)
	}
	return nil
}
