/*
Copyright 2020 The Crossplane Authors.

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

package eks

import (
	"testing"
	// "time"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/google/go-cmp/cmp"

	// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/provider-aws/apis/eks/v1alpha1"
)

var (
	fpName              = "my-cool-fargate-profile"
	podExecutionRoleArn = "arn:aws:iam::123456789:role/podexecutionrole"
	subnets             = []string{"subnet1", "subnet2"}
	namespace           = "fargate-namespace"
)

func TestGenerateCreateFargateProfileInput(t *testing.T) {
	type args struct {
		name string
		p    *v1alpha1.FargateProfileParameters
	}

	cases := map[string]struct {
		args args
		want *eks.CreateFargateProfileInput
	}{
		"AllFields": {
			args: args{
				name: fpName,
				p: &v1alpha1.FargateProfileParameters{
					ClusterName:         clusterName,
					PodExecutionRoleArn: podExecutionRoleArn,
					Subnets:             subnets,
					Tags:                map[string]string{"cool": "tag"},
					Selectors: []v1alpha1.FargateProfileSelector{
						{
							Namespace: &namespace,
							Labels: map[string]string{
								"cool": "label",
							},
						},
					},
				},
			},
			want: &eks.CreateFargateProfileInput{
				FargateProfileName:  &fpName,
				ClusterName:         &clusterName,
				PodExecutionRoleArn: &podExecutionRoleArn,
				Subnets:             subnets,
				Tags:                map[string]string{"cool": "tag"},
				Selectors: []eks.FargateProfileSelector{
					{
						Namespace: &namespace,
						Labels: map[string]string{
							"cool": "label",
						},
					},
				},
			},
		},
		"SomeFields": {
			args: args{
				name: fpName,
				p: &v1alpha1.FargateProfileParameters{
					ClusterName:         clusterName,
					PodExecutionRoleArn: podExecutionRoleArn,
					Selectors: []v1alpha1.FargateProfileSelector{
						{
							Namespace: &namespace,
							Labels: map[string]string{
								"cool": "label",
							},
						},
					},
				},
			},
			want: &eks.CreateFargateProfileInput{
				FargateProfileName:  &fpName,
				ClusterName:         &clusterName,
				PodExecutionRoleArn: &podExecutionRoleArn,
				Selectors: []eks.FargateProfileSelector{
					{
						Namespace: &namespace,
						Labels: map[string]string{
							"cool": "label",
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateCreateFargateProfileInput(tc.args.name, tc.args.p)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestLateInitializeFargateProfile(t *testing.T) {
	type args struct {
		p *v1alpha1.FargateProfileParameters
		n *eks.FargateProfile
	}

	cases := map[string]struct {
		args args
		want *v1alpha1.FargateProfileParameters
	}{
		"AllFieldsEmpty": {
			args: args{
				p: &v1alpha1.FargateProfileParameters{},
				n: &eks.FargateProfile{
					Subnets: subnets,
					Tags:    map[string]string{"cool": "tag"},
				},
			},
			want: &v1alpha1.FargateProfileParameters{
				Subnets: subnets,
				Tags:    map[string]string{"cool": "tag"},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			LateInitializeFargateProfile(tc.args.p, tc.args.n)
			if diff := cmp.Diff(tc.want, tc.args.p); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsFargateProfileUpToDate(t *testing.T) {
	type args struct {
		p *v1alpha1.FargateProfileParameters
		n *eks.FargateProfile
	}

	cases := map[string]struct {
		args args
		want bool
	}{
		"UpToDate": {
			args: args{
				p: &v1alpha1.FargateProfileParameters{
					Tags: map[string]string{"cool": "tag"},
				},
				n: &eks.FargateProfile{
					Tags: map[string]string{"cool": "tag"},
				},
			},
			want: true,
		},
		"UpdateTags": {
			args: args{
				p: &v1alpha1.FargateProfileParameters{
					Tags: map[string]string{"cool": "tag", "another": "tag"},
				},
				n: &eks.FargateProfile{
					Tags: map[string]string{"cool": "tag"},
				},
			},
			want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			upToDate := IsFargateProfileUpToDate(tc.args.p, tc.args.n)
			if diff := cmp.Diff(tc.want, upToDate); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}