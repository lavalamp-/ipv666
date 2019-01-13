package config

func init() {
	InitConfig()
}

//TODO add tests for getting target network

//
//func TestConfiguration_SetTargetNetworkSets(t *testing.T) {
//	targetNetwork = nil
//	_, network, _ := net.ParseCIDR("2600::1/64")
//	SetTargetNetwork(network)
//	assert.NotNil(t, targetNetwork)
//}
//
//func TestConfiguration_SetTargetNetworkSetsValue(t *testing.T) {
//	targetNetwork = nil
//	_, network, _ := net.ParseCIDR("2600::1/64")
//	SetTargetNetwork(network)
//	expected := network.String()
//	assert.EqualValues(t, expected, targetNetwork.String())
//}
//
//func TestConfiguration_GetTargetNetwork(t *testing.T) {
//	targetNetwork = nil
//	_, err := GetTargetNetwork()
//	assert.Nil(t, err)
//}
//
//func TestConfiguration_GetTargetNetworkNilCreates(t *testing.T) {
//	targetNetwork = nil
//	network, _ := GetTargetNetwork()
//	expected := "2000::/4"
//	assert.EqualValues(t, expected, network.String())
//}
//
//func TestConfiguration_GetTargetNetworkNotNilReturns(t *testing.T) {
//	targetNetwork = nil
//	_, network, _ := net.ParseCIDR("2600::1/64")
//	SetTargetNetwork(network)
//	expected := network.String()
//	newNet, _ := GetTargetNetwork()
//	assert.EqualValues(t, expected, newNet.String())
//}
