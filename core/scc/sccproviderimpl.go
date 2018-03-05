/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scc

import (
	"fmt"

	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/core/common/sysccprovider"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/peer"
)

// ProviderFactory implements the sysccprovider.SystemChaincodeProviderFactory
// interface and returns instances of sysccprovider.SystemChaincodeProvider
type ProviderFactory struct {
	Peer        peer.Operations
	PeerSupport peer.Support
}

// NewSystemChaincodeProvider returns pointers to ProviderFactory as an
// implementer of the sysccprovider.SystemChaincodeProvider interface
func (c *ProviderFactory) NewSystemChaincodeProvider() sysccprovider.SystemChaincodeProvider {
	return &sccProviderImpl{
		Peer:        c.Peer,
		PeerSupport: c.PeerSupport,
	}
}

// init is called when this package is loaded. This implementation registers the factory
func init() {
	sysccprovider.RegisterSystemChaincodeProviderFactory(
		&ProviderFactory{
			Peer:        peer.Default,
			PeerSupport: peer.DefaultSupport,
		},
	)
}

// ccProviderImpl is an implementation of the ccprovider.ChaincodeProvider interface
type sccProviderImpl struct {
	Peer        peer.Operations
	PeerSupport peer.Support
}

// IsSysCC returns true if the supplied chaincode is a system chaincode
func (c *sccProviderImpl) IsSysCC(name string) bool {
	return IsSysCC(name)
}

// IsSysCCAndNotInvokableCC2CC returns true if the supplied chaincode is
// ia system chaincode and it NOT invokable through a cc2cc invocation
func (c *sccProviderImpl) IsSysCCAndNotInvokableCC2CC(name string) bool {
	return IsSysCCAndNotInvokableCC2CC(name)
}

// GetQueryExecutorForLedger returns a query executor for the specified channel
func (c *sccProviderImpl) GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error) {
	l := c.Peer.GetLedger(cid)
	if l == nil {
		return nil, fmt.Errorf("Could not retrieve ledger for channel %s", cid)
	}

	return l.NewQueryExecutor()
}

// IsSysCCAndNotInvokableExternal returns true if the supplied chaincode is
// ia system chaincode and it NOT invokable
func (c *sccProviderImpl) IsSysCCAndNotInvokableExternal(name string) bool {
	// call the static method of the same name
	return IsSysCCAndNotInvokableExternal(name)
}

// GetApplicationConfig returns the configtxapplication.SharedConfig for the channel
// and whether the Application config exists
func (c *sccProviderImpl) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return c.PeerSupport.GetApplicationConfig(cid)
}

// Returns the policy manager associated to the passed channel
// and whether the policy manager exists
func (c *sccProviderImpl) PolicyManager(channelID string) (policies.Manager, bool) {
	m := c.Peer.GetPolicyManager(channelID)
	return m, (m != nil)
}
