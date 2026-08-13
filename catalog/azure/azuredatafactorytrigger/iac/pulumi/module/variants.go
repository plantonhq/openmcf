package module

import (
	"github.com/pkg/errors"
	azuredatafactorytriggerv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorytrigger/v1alpha1"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/datafactory"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Each variant creates its matching provider resource. The SDK
// generates a distinct pipeline-reference type per resource, so each
// variant builds its own list -- all write the same
// PipelineReference wire shape.

func createScheduleTrigger(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	schedule := spec.Schedule

	pipelines := datafactory.TriggerSchedulePipelineArray{}
	for _, pipelineRef := range schedule.Pipelines {
		pipelineArgs := datafactory.TriggerSchedulePipelineArgs{
			Name: pulumi.String(pipelineRef.Name.GetValue()),
		}
		if len(pipelineRef.Parameters) > 0 {
			pipelineArgs.Parameters = pulumi.ToStringMap(pipelineRef.Parameters)
		}
		pipelines = append(pipelines, pipelineArgs)
	}

	args := &datafactory.TriggerScheduleArgs{
		Name:          pulumi.String(spec.Name),
		DataFactoryId: pulumi.String(spec.DataFactoryId.GetValue()),
		Activated:     pulumi.Bool(activatedOrDefault(spec)),
		// The platform defaults, sent explicitly (the provider's own
		// defaults, mirrored).
		Frequency: pulumi.String(frequencyOrDefault(schedule)),
		Interval:  pulumi.Int(intervalOrDefault(schedule)),
		Pipelines: pipelines,
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringArray(spec.Annotations)
	}
	// Omitted start_time lets Azure fill in the moment of deployment
	// (the provider injects now-UTC itself when unset).
	if schedule.StartTime != "" {
		args.StartTime = pulumi.String(schedule.StartTime)
	}
	if schedule.EndTime != "" {
		args.EndTime = pulumi.String(schedule.EndTime)
	}
	if schedule.TimeZone != "" {
		args.TimeZone = pulumi.String(schedule.TimeZone)
	}

	if recurrence := schedule.RecurrenceSchedule; recurrence != nil {
		// ENGINE-SHAPE: the bridged SDK pluralizes the provider's
		// list-arg names (DaysOfMonths, DaysOfWeeks, Monthlies) -- name
		// differences only; both engines write the same ARM recurrence
		// schedule.
		scheduleArgs := datafactory.TriggerScheduleScheduleArgs{}
		if len(recurrence.DaysOfMonth) > 0 {
			daysOfMonth := make([]int, 0, len(recurrence.DaysOfMonth))
			for _, day := range recurrence.DaysOfMonth {
				daysOfMonth = append(daysOfMonth, int(day))
			}
			scheduleArgs.DaysOfMonths = pulumi.ToIntArray(daysOfMonth)
		}
		if len(recurrence.DaysOfWeek) > 0 {
			scheduleArgs.DaysOfWeeks = pulumi.ToStringArray(recurrence.DaysOfWeek)
		}
		if len(recurrence.Hours) > 0 {
			hours := make([]int, 0, len(recurrence.Hours))
			for _, hour := range recurrence.Hours {
				hours = append(hours, int(hour))
			}
			scheduleArgs.Hours = pulumi.ToIntArray(hours)
		}
		if len(recurrence.Minutes) > 0 {
			minutes := make([]int, 0, len(recurrence.Minutes))
			for _, minute := range recurrence.Minutes {
				minutes = append(minutes, int(minute))
			}
			scheduleArgs.Minutes = pulumi.ToIntArray(minutes)
		}
		if len(recurrence.Monthly) > 0 {
			monthlies := datafactory.TriggerScheduleScheduleMonthlyArray{}
			for _, occurrence := range recurrence.Monthly {
				monthlyArgs := datafactory.TriggerScheduleScheduleMonthlyArgs{
					Weekday: pulumi.String(occurrence.Weekday),
				}
				if occurrence.Week != nil {
					monthlyArgs.Week = pulumi.IntPtr(int(occurrence.GetWeek()))
				}
				monthlies = append(monthlies, monthlyArgs)
			}
			scheduleArgs.Monthlies = monthlies
		}
		args.Schedule = scheduleArgs
	}

	createdTrigger, err := datafactory.NewTriggerSchedule(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create data factory schedule trigger %s", resourceName)
	}
	return createdTrigger.ID(), createdTrigger.Name, nil
}

func frequencyOrDefault(schedule *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSchedule) string {
	if schedule.Frequency != nil && schedule.GetFrequency() != "" {
		return schedule.GetFrequency()
	}
	return "Minute"
}

func intervalOrDefault(schedule *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSchedule) int {
	if schedule.Interval != nil {
		return int(schedule.GetInterval())
	}
	return 1
}

func createTumblingWindowTrigger(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	tumblingWindow := spec.TumblingWindow

	// Tumbling window triggers drive exactly ONE pipeline (Azure's own
	// model).
	pipelineArgs := datafactory.TriggerTumblingWindowPipelineArgs{
		Name: pulumi.String(tumblingWindow.Pipeline.Name.GetValue()),
	}
	if len(tumblingWindow.Pipeline.Parameters) > 0 {
		pipelineArgs.Parameters = pulumi.ToStringMap(tumblingWindow.Pipeline.Parameters)
	}

	maxConcurrency := 50
	if tumblingWindow.MaxConcurrency != nil {
		maxConcurrency = int(tumblingWindow.GetMaxConcurrency())
	}

	args := &datafactory.TriggerTumblingWindowArgs{
		Name:           pulumi.String(spec.Name),
		DataFactoryId:  pulumi.String(spec.DataFactoryId.GetValue()),
		Activated:      pulumi.Bool(activatedOrDefault(spec)),
		Frequency:      pulumi.String(tumblingWindow.Frequency),
		Interval:       pulumi.Int(int(tumblingWindow.Interval)),
		StartTime:      pulumi.String(tumblingWindow.StartTime),
		MaxConcurrency: pulumi.Int(maxConcurrency),
		Pipeline:       pipelineArgs,
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringArray(spec.Annotations)
	}
	if tumblingWindow.EndTime != "" {
		args.EndTime = pulumi.String(tumblingWindow.EndTime)
	}
	if tumblingWindow.Delay != "" {
		args.Delay = pulumi.String(tumblingWindow.Delay)
	}
	if len(tumblingWindow.AdditionalProperties) > 0 {
		args.AdditionalProperties = pulumi.ToStringMap(tumblingWindow.AdditionalProperties)
	}

	if tumblingWindow.Retry != nil {
		retryInterval := 30
		if tumblingWindow.Retry.Interval != nil {
			retryInterval = int(tumblingWindow.Retry.GetInterval())
		}
		args.Retry = datafactory.TriggerTumblingWindowRetryArgs{
			Count:    pulumi.Int(int(tumblingWindow.Retry.Count)),
			Interval: pulumi.IntPtr(retryInterval),
		}
	}

	if len(tumblingWindow.Dependencies) > 0 {
		dependencies := datafactory.TriggerTumblingWindowTriggerDependencyArray{}
		for _, dependency := range tumblingWindow.Dependencies {
			dependencyArgs := datafactory.TriggerTumblingWindowTriggerDependencyArgs{}
			// An entry without trigger_name is a SELF-dependency (this
			// trigger's own earlier windows) -- the provider's own
			// convention.
			if dependency.TriggerName != nil && dependency.TriggerName.GetValue() != "" {
				dependencyArgs.TriggerName = pulumi.String(dependency.TriggerName.GetValue())
			}
			if dependency.Offset != "" {
				dependencyArgs.Offset = pulumi.String(dependency.Offset)
			}
			if dependency.Size != "" {
				dependencyArgs.Size = pulumi.String(dependency.Size)
			}
			dependencies = append(dependencies, dependencyArgs)
		}
		args.TriggerDependencies = dependencies
	}

	createdTrigger, err := datafactory.NewTriggerTumblingWindow(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create data factory tumbling window trigger %s", resourceName)
	}
	return createdTrigger.ID(), createdTrigger.Name, nil
}

func createBlobEventTrigger(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	blobEvent := spec.BlobEvent

	pipelines := datafactory.TriggerBlobEventPipelineArray{}
	for _, pipelineRef := range blobEvent.Pipelines {
		pipelineArgs := datafactory.TriggerBlobEventPipelineArgs{
			Name: pulumi.String(pipelineRef.Name.GetValue()),
		}
		if len(pipelineRef.Parameters) > 0 {
			pipelineArgs.Parameters = pulumi.ToStringMap(pipelineRef.Parameters)
		}
		pipelines = append(pipelines, pipelineArgs)
	}

	args := &datafactory.TriggerBlobEventArgs{
		Name:             pulumi.String(spec.Name),
		DataFactoryId:    pulumi.String(spec.DataFactoryId.GetValue()),
		Activated:        pulumi.Bool(activatedOrDefault(spec)),
		StorageAccountId: pulumi.String(blobEvent.StorageAccountId.GetValue()),
		Events:           pulumi.ToStringArray(blobEvent.Events),
		// The platform default, sent explicitly (the provider sends the
		// effective bool unconditionally too).
		IgnoreEmptyBlobs: pulumi.Bool(ignoreEmptyBlobsOrDefault(blobEvent)),
		Pipelines:        pipelines,
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringArray(spec.Annotations)
	}
	if blobEvent.BlobPathBeginsWith != "" {
		args.BlobPathBeginsWith = pulumi.String(blobEvent.BlobPathBeginsWith)
	}
	if blobEvent.BlobPathEndsWith != "" {
		args.BlobPathEndsWith = pulumi.String(blobEvent.BlobPathEndsWith)
	}
	if len(blobEvent.AdditionalProperties) > 0 {
		args.AdditionalProperties = pulumi.ToStringMap(blobEvent.AdditionalProperties)
	}

	createdTrigger, err := datafactory.NewTriggerBlobEvent(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create data factory blob event trigger %s", resourceName)
	}
	return createdTrigger.ID(), createdTrigger.Name, nil
}

func ignoreEmptyBlobsOrDefault(blobEvent *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerBlobEvent) bool {
	if blobEvent.IgnoreEmptyBlobs != nil {
		return blobEvent.GetIgnoreEmptyBlobs()
	}
	return false
}

func createCustomEventTrigger(
	ctx *pulumi.Context,
	resourceName string,
	spec *azuredatafactorytriggerv1alpha1.AzureDataFactoryTriggerSpec,
	azureProvider pulumi.ProviderResource,
) (pulumi.StringInput, pulumi.StringInput, error) {
	customEvent := spec.CustomEvent

	pipelines := datafactory.TriggerCustomEventPipelineArray{}
	for _, pipelineRef := range customEvent.Pipelines {
		pipelineArgs := datafactory.TriggerCustomEventPipelineArgs{
			Name: pulumi.String(pipelineRef.Name.GetValue()),
		}
		if len(pipelineRef.Parameters) > 0 {
			pipelineArgs.Parameters = pulumi.ToStringMap(pipelineRef.Parameters)
		}
		pipelines = append(pipelines, pipelineArgs)
	}

	args := &datafactory.TriggerCustomEventArgs{
		Name:             pulumi.String(spec.Name),
		DataFactoryId:    pulumi.String(spec.DataFactoryId.GetValue()),
		Activated:        pulumi.Bool(activatedOrDefault(spec)),
		EventgridTopicId: pulumi.String(customEvent.EventgridTopicId.GetValue()),
		Events:           pulumi.ToStringArray(customEvent.Events),
		Pipelines:        pipelines,
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringArray(spec.Annotations)
	}
	if customEvent.SubjectBeginsWith != "" {
		args.SubjectBeginsWith = pulumi.String(customEvent.SubjectBeginsWith)
	}
	if customEvent.SubjectEndsWith != "" {
		args.SubjectEndsWith = pulumi.String(customEvent.SubjectEndsWith)
	}
	if len(customEvent.AdditionalProperties) > 0 {
		args.AdditionalProperties = pulumi.ToStringMap(customEvent.AdditionalProperties)
	}

	createdTrigger, err := datafactory.NewTriggerCustomEvent(ctx, resourceName, args, pulumi.Provider(azureProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create data factory custom event trigger %s", resourceName)
	}
	return createdTrigger.ID(), createdTrigger.Name, nil
}
