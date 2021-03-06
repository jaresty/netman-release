// This file was generated by counterfeiter
package fakes

import (
	"net"
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
)

type NetOutProvider struct {
	InitializeStub        func(logger lager.Logger, containerHandle string, containerIP net.IP, overlayNetwork string) error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		logger          lager.Logger
		containerHandle string
		containerIP     net.IP
		overlayNetwork  string
	}
	initializeReturns struct {
		result1 error
	}
	CleanupStub        func(containerHandle string) error
	cleanupMutex       sync.RWMutex
	cleanupArgsForCall []struct {
		containerHandle string
	}
	cleanupReturns struct {
		result1 error
	}
	InsertRuleStub        func(containerHandle string, rule garden.NetOutRule, containerIP string) error
	insertRuleMutex       sync.RWMutex
	insertRuleArgsForCall []struct {
		containerHandle string
		rule            garden.NetOutRule
		containerIP     string
	}
	insertRuleReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *NetOutProvider) Initialize(logger lager.Logger, containerHandle string, containerIP net.IP, overlayNetwork string) error {
	fake.initializeMutex.Lock()
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
		logger          lager.Logger
		containerHandle string
		containerIP     net.IP
		overlayNetwork  string
	}{logger, containerHandle, containerIP, overlayNetwork})
	fake.recordInvocation("Initialize", []interface{}{logger, containerHandle, containerIP, overlayNetwork})
	fake.initializeMutex.Unlock()
	if fake.InitializeStub != nil {
		return fake.InitializeStub(logger, containerHandle, containerIP, overlayNetwork)
	} else {
		return fake.initializeReturns.result1
	}
}

func (fake *NetOutProvider) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *NetOutProvider) InitializeArgsForCall(i int) (lager.Logger, string, net.IP, string) {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return fake.initializeArgsForCall[i].logger, fake.initializeArgsForCall[i].containerHandle, fake.initializeArgsForCall[i].containerIP, fake.initializeArgsForCall[i].overlayNetwork
}

func (fake *NetOutProvider) InitializeReturns(result1 error) {
	fake.InitializeStub = nil
	fake.initializeReturns = struct {
		result1 error
	}{result1}
}

func (fake *NetOutProvider) Cleanup(containerHandle string) error {
	fake.cleanupMutex.Lock()
	fake.cleanupArgsForCall = append(fake.cleanupArgsForCall, struct {
		containerHandle string
	}{containerHandle})
	fake.recordInvocation("Cleanup", []interface{}{containerHandle})
	fake.cleanupMutex.Unlock()
	if fake.CleanupStub != nil {
		return fake.CleanupStub(containerHandle)
	} else {
		return fake.cleanupReturns.result1
	}
}

func (fake *NetOutProvider) CleanupCallCount() int {
	fake.cleanupMutex.RLock()
	defer fake.cleanupMutex.RUnlock()
	return len(fake.cleanupArgsForCall)
}

func (fake *NetOutProvider) CleanupArgsForCall(i int) string {
	fake.cleanupMutex.RLock()
	defer fake.cleanupMutex.RUnlock()
	return fake.cleanupArgsForCall[i].containerHandle
}

func (fake *NetOutProvider) CleanupReturns(result1 error) {
	fake.CleanupStub = nil
	fake.cleanupReturns = struct {
		result1 error
	}{result1}
}

func (fake *NetOutProvider) InsertRule(containerHandle string, rule garden.NetOutRule, containerIP string) error {
	fake.insertRuleMutex.Lock()
	fake.insertRuleArgsForCall = append(fake.insertRuleArgsForCall, struct {
		containerHandle string
		rule            garden.NetOutRule
		containerIP     string
	}{containerHandle, rule, containerIP})
	fake.recordInvocation("InsertRule", []interface{}{containerHandle, rule, containerIP})
	fake.insertRuleMutex.Unlock()
	if fake.InsertRuleStub != nil {
		return fake.InsertRuleStub(containerHandle, rule, containerIP)
	} else {
		return fake.insertRuleReturns.result1
	}
}

func (fake *NetOutProvider) InsertRuleCallCount() int {
	fake.insertRuleMutex.RLock()
	defer fake.insertRuleMutex.RUnlock()
	return len(fake.insertRuleArgsForCall)
}

func (fake *NetOutProvider) InsertRuleArgsForCall(i int) (string, garden.NetOutRule, string) {
	fake.insertRuleMutex.RLock()
	defer fake.insertRuleMutex.RUnlock()
	return fake.insertRuleArgsForCall[i].containerHandle, fake.insertRuleArgsForCall[i].rule, fake.insertRuleArgsForCall[i].containerIP
}

func (fake *NetOutProvider) InsertRuleReturns(result1 error) {
	fake.InsertRuleStub = nil
	fake.insertRuleReturns = struct {
		result1 error
	}{result1}
}

func (fake *NetOutProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	fake.cleanupMutex.RLock()
	defer fake.cleanupMutex.RUnlock()
	fake.insertRuleMutex.RLock()
	defer fake.insertRuleMutex.RUnlock()
	return fake.invocations
}

func (fake *NetOutProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
