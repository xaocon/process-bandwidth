package main

import (
	"github.com/gojue/ebpfmanager"
	"github.com/sirupsen/logrus"
)

var m = &manager.Manager{
	Probes: []*manager.Probe{
		&manager.Probe{
			UID:              "MyFirstHook",
			Section:          "kprobe/vfs_mkdir",
			AttachToFuncName: "vfs_mkdir",
			EbpfFuncName:     "kprobe_vfs_mkdir",
		},
		&manager.Probe{
			UID:              "", // UID is not needed if there will be only one instance of the program
			Section:          "kretprobe/mkdir",
			AttachToFuncName: "mkdir",
			EbpfFuncName:     "kretpobe_mkdir",
			KProbeMaxActive:  100,
		},
	},
	PerfMaps: []*manager.PerfMap{
		&manager.PerfMap{
			Map: manager.Map{
				Name: "my_constants",
			},
			PerfMapOptions: manager.PerfMapOptions{
				DataHandler: myDataHandler,
			},
		},
	},
}

// myDataHandler - Perf event data handler
func myDataHandler(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
	myConstant := ByteOrder.Uint64(data[0:8])
	logrus.Printf("received: CPU:%d my_constant:%d", cpu, myConstant)
}

var editors = []manager.ConstantEditor{
	{
		Name:          "my_constant",
		Value:         uint64(100),
		FailOnMissing: true,
		ProbeIdentificationPairs: []manager.ProbeIdentificationPair{
			{UID: "MyFirstHook", EbpfFuncName: "kprobe_vfs_mkdir"},
		},
	},
	{
		Name:          "my_constant",
		Value:         uint64(555),
		FailOnMissing: true,
		ProbeIdentificationPairs: []manager.ProbeIdentificationPair{
			{UID: "", EbpfFuncName: "kprobe_vfs_rmdir"},
		},
	},
	{
		Name:                     "unused_constant",
		Value:                    uint64(555),
		ProbeIdentificationPairs: []manager.ProbeIdentificationPair{},
	},
}

func main() {
	// Prepare manager options
	options := manager.Options{ConstantEditors: editors}

	// Initialize the manager
	if err := m.InitWithOptions(recoverAssets(), options); err != nil {
		logrus.Fatal(err)
	}

	// Start the manager
	if err := m.Start(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Println("eBPF programs running, head over to /sys/kernel/debug/tracing/trace_pipe to see them in action.")

	// Demo
	logrus.Println("INITIAL PROGRAMS")
	if err := trigger(); err != nil {
		_ = m.Stop(manager.CleanAll)
		logrus.Fatal(err)
	}
	if err := demoClone(); err != nil {
		_ = m.Stop(manager.CleanAll)
		logrus.Fatal(err)
	}
	if err := demoAddHook(); err != nil {
		_ = m.Stop(manager.CleanAll)
		logrus.Fatal(err)
	}

	// Close the manager
	if err := m.Stop(manager.CleanAll); err != nil {
		logrus.Fatal(err)
	}
}
