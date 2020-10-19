// Copyright 2020 The Alibaba Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubeflow/arena/pkg/util"
	"github.com/kubeflow/arena/pkg/workflow"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	elasticdljobChart = util.GetChartsFolder() + "/elasticdljob"
)

func NewSubmitElasticDLJobCommand() *cobra.Command {
	var (
		submitArgs submitElasticDLJobArgs
	)

	submitArgs.Mode = "elasticdljob"

	var command = &cobra.Command{
		Use:     "elasticdljob",
		Short:   "Sumit ElasticDLJob as training job.",
		Aliases: []string{"elasticdl"},
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("hello world")
			log.Info(args)
			// if len(args) == 0 {
			// cmd.HelpFunc()(cmd, args)
			// os.Exit(1)
			// }

			util.SetLogLevel(logLevel)
			setupKubeconfig()
			_, err := initKubeClient()
			if err != nil {
				log.Debugf("Failed due to %v", err)
				fmt.Println(err)
				os.Exit(1)
			}

			err = updateNamespace(cmd)
			if err != nil {
				log.Debugf("Failed due to %v", err)
				fmt.Println(err)
				os.Exit(1)
			}

			err = submitElasticDLJob(args, &submitArgs)
			if err != nil {
				log.Debugf("Failed due to %v", err)
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	log.Debugf("elasticdljob command: %v", command)

	command.Flags().StringVar(&submitArgs.ModelZoo, "modelZoo", "", "")
	command.Flags().StringVar(&submitArgs.ModelDef, "modelDef", "", "")
	command.Flags().StringVar(&submitArgs.TrainingData, "trainingData", "", "")
	command.Flags().StringVar(&submitArgs.ValidationData, "validationData", "", "")
	command.Flags().StringVar(&submitArgs.Output, "output", "", "")
	command.Flags().IntVar(&submitArgs.NumEpochs, "numEpochs", 0, "")
	command.Flags().IntVar(&submitArgs.MinibatchSize, "minibatchSize", 0, "")
	command.Flags().IntVar(&submitArgs.NumMinibatchesPerTask, "numMinibatchesPerTask", 0, "")
	command.Flags().IntVar(&submitArgs.EvaluationSteps, "evaluationSteps", 0, "")

	command.Flags().StringVar(&submitArgs.ImagePullPolicy, "imagePullPolicy", "", "")
	command.Flags().StringVar(&submitArgs.Volume, "volume", "", "")

	command.Flags().StringVar(&submitArgs.MasterPriority, "masterPriority", "", "")
	command.Flags().StringVar(&submitArgs.MasterCPU, "masterCPU", "", "")
	command.Flags().StringVar(&submitArgs.MasterMemory, "masterMemory", "", "")

	command.Flags().IntVar(&submitArgs.PSCount, "psCount", 0, "")
	command.Flags().StringVar(&submitArgs.PSPriority, "psPriority", "", "")
	command.Flags().StringVar(&submitArgs.PSCPU, "psCPU", "", "")
	command.Flags().StringVar(&submitArgs.PSMemory, "psMemory", "", "")

	command.Flags().IntVar(&submitArgs.WorkerCount, "workerCount", 0, "")
	command.Flags().StringVar(&submitArgs.WorkerPriority, "workerPriority", "", "")
	command.Flags().StringVar(&submitArgs.WorkerCPU, "workerCPU", "", "")
	command.Flags().StringVar(&submitArgs.WorkerMemory, "workerMemory", "", "")

	submitArgs.addCommonFlags(command)
	log.Debugf("elasticdljob command: %v", command)
	return command
}

type submitElasticDLJobArgs struct {
	submitArgs            `yaml:",inline"`
	ModelZoo              string `yaml:"modelZoo"`
	ModelDef              string `yaml:"modelDef"`
	TrainingData          string `yaml:"trainingData"`
	ValidationData        string `yaml:"validationData"`
	Output                string `yaml:"output"`
	NumEpochs             int    `yaml:"numEpochs"`
	MinibatchSize         int    `yaml:"minibatchSize"`
	NumMinibatchesPerTask int    `yaml:"numMimibatchesPerTask"`
	EvaluationSteps       int    `yaml:"evaluationSteps"`
	ImagePullPolicy       string `yaml:"imagePullPolicy"`
	Volume                string `yaml:"volume"`
	MasterPriority        string `yaml:"masterPriority"`
	MasterCPU             string `yaml:"masterCPU"`
	MasterMemory          string `yaml:"masterMemory"`
	PSCount               int    `yaml:"psCount"`
	PSPriority            string `yaml:"psPriority"`
	PSCPU                 string `yaml:"psCPU"`
	PSMemory              string `yaml:"psMemory"`
	WorkerCount           int    `yaml:"workerCount"`
	WorkerPriority        string `yaml:"workerPriority"`
	WorkerCPU             string `yaml:"workerCPU"`
	WorkerMemory          string `yaml:"workerMemory"`
}

func (submitArgs *submitElasticDLJobArgs) prepare(args []string) (err error) {
	submitArgs.Command = strings.Join(args, " ")

	commonArgs := &submitArgs.submitArgs

	err = commonArgs.transform()
	if err != nil {
		return err
	}

	if err := submitArgs.addConfigFiles(); err != nil {
		return err
	}

	if len(envs) > 0 {
		submitArgs.Envs = transformSliceToMap(envs, "=")
	}

	submitArgs.processCommonFlags()

	return nil
}

func (submitArgs *submitElasticDLJobArgs) addConfigFiles() error {
	return submitArgs.addJobConfigFiles()
}

func submitElasticDLJob(args []string, submitArgs *submitElasticDLJobArgs) (err error) {
	err = submitArgs.prepare(args)
	if err != nil {
		return err
	}

	err = workflow.SubmitJob(name, submitArgs.Mode, namespace, submitArgs, elasticdljobChart, submitArgs.addHelmOptions()...)
	if err != nil {
		return err
	}

	log.Infof("The Job %s has been submitted successfully", name)
	log.Infof("You can run `arena get %s --type %s` to check the job status", name, submitArgs.Mode)

	return nil
}
