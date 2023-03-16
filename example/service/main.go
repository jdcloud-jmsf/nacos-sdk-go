/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func main() {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("11.50.138.27", 31718,
			constant.WithContextPath("/nacos"),
			constant.WithGrpcPort(31030)),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithUsername("nacos"),
		constant.WithPassword("nacos"),
		constant.WithUpdateCacheWhenEmpty(true),
	)

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	getNamespaces(client)

	//Register
	registerServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//Register
	updateServiceInstance(client, vo.UpdateInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      false,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//DeRegister
	//deRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.10",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Cluster:     "cluster-a",
	//	Ephemeral:   true, //it must be true
	//})

	time.Sleep(1 * time.Second)

	////BatchRegister
	//batchRegisterServiceInstance(client, vo.BatchRegisterInstanceParam{
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Instances: []vo.RegisterInstanceParam{{
	//		Ip:          "10.0.0.10",
	//		Port:        8848,
	//		Weight:      10,
	//		Enable:      true,
	//		Healthy:     true,
	//		Ephemeral:   true,
	//		ClusterName: "cluster-a",
	//		Metadata:    map[string]string{"idc": "shanghai"},
	//	}, {
	//		Ip:          "10.0.0.12",
	//		Port:        8848,
	//		Weight:      7,
	//		Enable:      true,
	//		Healthy:     true,
	//		Ephemeral:   true,
	//		ClusterName: "cluster-a",
	//		Metadata:    map[string]string{"idc": "shanghai"},
	//	}},
	//})
	//
	//time.Sleep(1 * time.Second)

	//Get service with serviceName, groupName , clusters
	//getService(client, vo.GetServiceParam{
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Clusters:    []string{"cluster-a"},
	//})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	//selectAllInstances(client, vo.SelectAllInstancesParam{
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Clusters:    []string{"cluster-a"},
	//})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	//selectInstances(client, vo.SelectInstancesParam{
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Clusters:    []string{"cluster-a"},
	//})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	//selectOneHealthyInstance(client, vo.SelectOneHealthInstanceParam{
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Clusters:    []string{"cluster-a"},
	//})

	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	subscribeParam := &vo.SubscribeParam{
		ServiceName: "jmsf-provider-demo",
		GroupName:   "DEFAULT_GROUP",
		SubscribeCallback: func(service string, services []model.Instance, err error) {
			fmt.Printf("callback return services:%s \n\n", util.ToJsonString(services))
		},
	}
	subscribe(client, subscribeParam)

	getCatalogServices(client, "public")

	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	/*updateServiceInstance(client, vo.UpdateInstanceParam{
		Ip:          "10.0.0.11", //update ip
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "beijing1"}, //update metadata
	})*/

	//wait for client pull change from server
	time.Sleep(10 * time.Minute)
	// UnSubscribe
	unSubscribe(client, subscribeParam)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	getAllService(client, vo.GetAllServiceInfoParam{
		GroupName: "group-a",
		PageNo:    1,
		PageSize:  10,
	})
}
