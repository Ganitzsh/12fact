package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/ganitzsh/12fact/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type BlurCmdArgs struct {
	CmdArgs
	pb.BlurImageRequest
	File string
}

func NewBlurCmdArgs() *BlurCmdArgs {
	return &BlurCmdArgs{}
}

func (a *BlurCmdArgs) Read(cmd *cobra.Command, args []string) error {
	sigma, _ := strconv.ParseFloat(args[1], 32)
	a.Sigma = float32(sigma)
	a.File = args[0]
	return nil
}

func blurCmdRun(cmd *cobra.Command, args []string) {
	a := NewBlurCmdArgs()
	if err := a.Read(cmd, args); err != nil {
		logrus.Errorf("Could not read arguments: %v", err)
		os.Exit(1)
	}
	logrus.Info("Blurring ", a.File, " with a factor of ", a.Sigma)
	client, err := rpcv1.NewClient()
	if err != nil {
		logrus.Errorf("Failed to get client: %v", err)
		os.Exit(1)
	}
	r, err := client.Transform(
		a.File, pb.TransformationType_BLUR, &pb.BlurImageRequest{
			Sigma: a.Sigma,
		})
	if err != nil {
		logrus.Errorf("Failed to blur image: %v", err)
		os.Exit(1)
	}
	outPath := "./out" + filepath.Ext(a.File)
	if outFlag, _ := cmd.Flags().GetString("out"); outFlag != "" {
		outPath = outFlag
	}
	if err := client.SaveToFile(r, outPath); err != nil {
		logrus.Errorf("Failed to save to %v: %ev", outPath, err)
	}
}

// blurCmd represents the serve command
var blurCmd = &cobra.Command{
	Use:   "blur [file] [sigma]",
	Short: "Blurs the given image with a factor of [sigma]",
	Run:   blurCmdRun,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errNoArgs
		}
		if len(args) != 2 {
			return errMissingArgs
		}
		if _, err := strconv.ParseFloat(args[0], 32); err != nil {
			return errors.New("sigma must be a floating number")
		}
		return nil
	}),
}

func init() {
	rootCmd.AddCommand(blurCmd)

	blurCmd.PersistentFlags().StringP(
		"out", "o", "", "The path to the output file",
	)
}
