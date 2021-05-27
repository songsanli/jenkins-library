package cmd

import (
	"encoding/json"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteCPECommand() *cobra.Command {
	return &cobra.Command{
		Use:   "writeCPE",
		Short: "Serializes the commonPipelineEnvironment JSON to disk",
		PreRun: func(cmd *cobra.Command, args []string) {
			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)
		},

		Run: func(cmd *cobra.Command, args []string) {
			err := runWriteCPE()
			if err != nil {
				log.Entry().Fatalf("error when writing common Pipeline environment: %v", err)
			}
		},
	}
}

func runWriteCPE() error {
	inBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	commonPipelineEnv := piperenv.CPEMap{}
	err = json.Unmarshal(inBytes, &commonPipelineEnv)
	if err != nil {
	}

	commonPipelineEnv.Flatten()

	rootPath := filepath.Join(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
	err = commonPipelineEnv.WriteToDisk(rootPath)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(commonPipelineEnv, "", "\t")
	os.Stdout.Write(bytes)
	return nil
}
