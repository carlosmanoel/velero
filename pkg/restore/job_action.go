/*
Copyright the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restore

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	batchv1api "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	riav1 "github.com/vmware-tanzu/velero/pkg/plugin/velero/restoreitemaction/v1"
)

type JobAction struct {
	logger logrus.FieldLogger
}

func NewJobAction(logger logrus.FieldLogger) *JobAction {
	return &JobAction{logger: logger}
}

func (a *JobAction) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"jobs"},
	}, nil
}

func (a *JobAction) Execute(input *riav1.RestoreItemActionExecuteInput) (*riav1.RestoreItemActionExecuteOutput, error) {
	job := new(batchv1api.Job)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(input.Item.UnstructuredContent(), job); err != nil {
		return nil, errors.WithStack(err)
	}

	if job.Spec.Selector != nil {
		delete(job.Spec.Selector.MatchLabels, "controller-uid")
	}
	delete(job.Spec.Template.ObjectMeta.Labels, "controller-uid")

	res, err := runtime.DefaultUnstructuredConverter.ToUnstructured(job)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return riav1.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: res}), nil
}
