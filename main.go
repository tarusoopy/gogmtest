package main

import (
	"context"
	"fmt"
	"github.com/mindstand/gogm/v2"
	"time"
)

type tdString string
type tdInt int

//structs for the example (can also be found in decoder_test.go)
type VertexA struct {
	// provides required node fields
	gogm.BaseNode

	TestField         string            `gogm:"name=test_field"`
	TestTypeDefString tdString          `gogm:"name=test_type_def_string"`
	TestTypeDefInt    tdInt             `gogm:"name=test_type_def_int"`
	MapProperty       map[string]string `gogm:"name=map_property;properties"`
	SliceProperty     []string          `gogm:"name=slice_property;properties"`
	SingleA           *VertexB          `gogm:"direction=incoming;relationship=test_rel"`
	ManyA             []*VertexB        `gogm:"direction=incoming;relationship=testm2o"`
	MultiA            []*VertexB        `gogm:"direction=incoming;relationship=multib"`
	SingleSpecA       *EdgeC            `gogm:"direction=outgoing;relationship=special_single"`
	MultiSpecA        []*EdgeC          `gogm:"direction=outgoing;relationship=special_multi"`
}

type VertexB struct {
	// provides required node fields
	gogm.BaseNode

	TestField  string     `gogm:"name=test_field"`
	TestTime   time.Time  `gogm:"name=test_time"`
	Single     *VertexA   `gogm:"direction=outgoing;relationship=test_rel"`
	ManyB      *VertexA   `gogm:"direction=outgoing;relationship=testm2o"`
	Multi      []*VertexA `gogm:"direction=outgoing;relationship=multib"`
	SingleSpec *EdgeC     `gogm:"direction=incoming;relationship=special_single"`
	MultiSpec  []*EdgeC   `gogm:"direction=incoming;relationship=special_multi"`
}

type EdgeC struct {
	// provides required node fields
	gogm.BaseNode

	Start *VertexA
	End   *VertexB
	Test  string `gogm:"name=test"`
}

func main() {
	// define your configuration
	config := gogm.Config{
		Host:          "0.0.0.0",
		Port:          7687,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Username:      "neo4j",
		Password:      "hn3437",
		PoolSize:      50,
		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		TargetDbs:     nil,
		// default logger wraps the go "log" package, implement the Logger interface from gogm to use your own logger
		Logger: gogm.GetDefaultLogger(),
		// define the log level
		LogLevel: "DEBUG",
		// enable neo4j go driver to log
		EnableDriverLogs: false,
		// enable gogm to log params in cypher queries. WARNING THIS IS A SECURITY RISK! Only use this when debugging
		EnableLogParams: false,
		// enable open tracing. Ensure contexts have spans already. GoGM does not make root spans, only child spans
		OpentracingEnabled: false,
	}

	// register all vertices and edges
	// this is so that GoGM doesn't have to do reflect processing of each edge in real time
	// use nil or gogm.DefaultPrimaryKeyStrategy if you only want graph ids
	// we are using the default key strategy since our vertices are using BaseNode
	_gogm, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &VertexA{}, &VertexB{}, &EdgeC{})
	if err != nil {
		panic(err)
	}

	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := _gogm.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer sess.Close()

	aVal := &VertexA{
		TestField: "woo neo4j",
	}

	bVal := &VertexB{
		TestTime: time.Now().UTC(),
	}

	//set bi directional pointer
	bVal.Single = aVal
	aVal.SingleA = bVal

	err = sess.SaveDepth(context.Background(), aVal, 2)
	if err != nil {
		panic(err)
	}

	//load the object we just made (save will set the uuid)
	var readin VertexA
	err = sess.Load(context.Background(), &readin, aVal.Id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", readin)
}
