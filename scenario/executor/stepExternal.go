package scenexec

import (
	scenio "github.com/multiversx/mx-chain-scenario-go/scenario/io"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// ExecuteExternalStep executes an external step referenced by the scenario.
func (ae *ScenarioExecutor) ExecuteExternalStep(step *scenmodel.ExternalStepsStep) error {
	log.Trace("ExternalStepsStep", "path", step.Path)
	if len(step.Comment) > 0 {
		log.Trace("ExternalStepsStep", "comment", step.Comment)
	}

	fileResolverBackup := ae.fileResolver
	clonedFileResolver := ae.fileResolver.Clone()
	externalStepsRunner := scenio.NewScenarioController(ae, clonedFileResolver, ae.vmBuilder.GetVMType())

	extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
	setExternalStepGasTracing(ae, step)

	err := externalStepsRunner.RunSingleJSONScenario(extAbsPth, scenio.DefaultRunScenarioOptions())
	if err != nil {
		return err
	}

	ae.fileResolver = fileResolverBackup

	return nil
}
