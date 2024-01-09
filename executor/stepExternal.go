package scenexec

import (
	mc "github.com/multiversx/mx-chain-scenario-go/controller"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
)

// ExecuteExternalStep executes an external step referenced by the scenario.
func (ae *ScenarioExecutor) ExecuteExternalStep(step *mj.ExternalStepsStep) error {
	log.Trace("ExternalStepsStep", "path", step.Path)
	if len(step.Comment) > 0 {
		log.Trace("ExternalStepsStep", "comment", step.Comment)
	}

	fileResolverBackup := ae.fileResolver
	clonedFileResolver := ae.fileResolver.Clone()
	externalStepsRunner := mc.NewScenarioController(ae, clonedFileResolver)

	extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
	setExternalStepGasTracing(ae, step)

	err := externalStepsRunner.RunSingleJSONScenario(extAbsPth, mc.DefaultRunScenarioOptions())
	if err != nil {
		return err
	}

	ae.fileResolver = fileResolverBackup

	return nil
}
