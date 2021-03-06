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

// Code generated by ack-generate. DO NOT EDIT.

package table

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/dynamodb"
	svcsdkapi "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane/provider-aws/apis/dynamodb/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an Table resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create Table in AWS"
	errDescribe      = "failed to describe Table"
	errDelete        = "failed to delete Table"
)

type connector struct {
	kube client.Client
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.Table)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := awsclient.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, err
	}
	return &external{client: svcapi.New(sess), kube: c.kube}, errors.Wrap(err, errCreateSession)
}

type external struct {
	kube   client.Client
	client svcsdkapi.DynamoDBAPI
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.Table)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if err := e.preObserve(ctx, cr); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeTableInput(cr)
	resp, err := e.client.DescribeTableWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, errors.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	lateInitialize(&cr.Spec.ForProvider, resp)
	GenerateTable(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        isUpToDate(cr, resp),
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.Table)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	if err := e.preCreate(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	input := GenerateCreateTableInput(cr)
	resp, err := e.client.CreateTableWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}

	if resp.TableDescription.ArchivalSummary != nil {
		f0 := &svcapitypes.ArchivalSummary{}
		if resp.TableDescription.ArchivalSummary.ArchivalBackupArn != nil {
			f0.ArchivalBackupARN = resp.TableDescription.ArchivalSummary.ArchivalBackupArn
		}
		if resp.TableDescription.ArchivalSummary.ArchivalDateTime != nil {
			f0.ArchivalDateTime = &metav1.Time{*resp.TableDescription.ArchivalSummary.ArchivalDateTime}
		}
		if resp.TableDescription.ArchivalSummary.ArchivalReason != nil {
			f0.ArchivalReason = resp.TableDescription.ArchivalSummary.ArchivalReason
		}
		cr.Status.AtProvider.ArchivalSummary = f0
	}
	if resp.TableDescription.BillingModeSummary != nil {
		f2 := &svcapitypes.BillingModeSummary{}
		if resp.TableDescription.BillingModeSummary.BillingMode != nil {
			f2.BillingMode = resp.TableDescription.BillingModeSummary.BillingMode
		}
		if resp.TableDescription.BillingModeSummary.LastUpdateToPayPerRequestDateTime != nil {
			f2.LastUpdateToPayPerRequestDateTime = &metav1.Time{*resp.TableDescription.BillingModeSummary.LastUpdateToPayPerRequestDateTime}
		}
		cr.Status.AtProvider.BillingModeSummary = f2
	}
	if resp.TableDescription.CreationDateTime != nil {
		cr.Status.AtProvider.CreationDateTime = &metav1.Time{*resp.TableDescription.CreationDateTime}
	}
	if resp.TableDescription.GlobalTableVersion != nil {
		cr.Status.AtProvider.GlobalTableVersion = resp.TableDescription.GlobalTableVersion
	}
	if resp.TableDescription.ItemCount != nil {
		cr.Status.AtProvider.ItemCount = resp.TableDescription.ItemCount
	}
	if resp.TableDescription.LatestStreamArn != nil {
		cr.Status.AtProvider.LatestStreamARN = resp.TableDescription.LatestStreamArn
	}
	if resp.TableDescription.LatestStreamLabel != nil {
		cr.Status.AtProvider.LatestStreamLabel = resp.TableDescription.LatestStreamLabel
	}
	if resp.TableDescription.Replicas != nil {
		f12 := []*svcapitypes.ReplicaDescription{}
		for _, f12iter := range resp.TableDescription.Replicas {
			f12elem := &svcapitypes.ReplicaDescription{}
			if f12iter.GlobalSecondaryIndexes != nil {
				f12elemf0 := []*svcapitypes.ReplicaGlobalSecondaryIndexDescription{}
				for _, f12elemf0iter := range f12iter.GlobalSecondaryIndexes {
					f12elemf0elem := &svcapitypes.ReplicaGlobalSecondaryIndexDescription{}
					if f12elemf0iter.IndexName != nil {
						f12elemf0elem.IndexName = f12elemf0iter.IndexName
					}
					if f12elemf0iter.ProvisionedThroughputOverride != nil {
						f12elemf0elemf1 := &svcapitypes.ProvisionedThroughputOverride{}
						if f12elemf0iter.ProvisionedThroughputOverride.ReadCapacityUnits != nil {
							f12elemf0elemf1.ReadCapacityUnits = f12elemf0iter.ProvisionedThroughputOverride.ReadCapacityUnits
						}
						f12elemf0elem.ProvisionedThroughputOverride = f12elemf0elemf1
					}
					f12elemf0 = append(f12elemf0, f12elemf0elem)
				}
				f12elem.GlobalSecondaryIndexes = f12elemf0
			}
			if f12iter.KMSMasterKeyId != nil {
				f12elem.KMSMasterKeyID = f12iter.KMSMasterKeyId
			}
			if f12iter.ProvisionedThroughputOverride != nil {
				f12elemf2 := &svcapitypes.ProvisionedThroughputOverride{}
				if f12iter.ProvisionedThroughputOverride.ReadCapacityUnits != nil {
					f12elemf2.ReadCapacityUnits = f12iter.ProvisionedThroughputOverride.ReadCapacityUnits
				}
				f12elem.ProvisionedThroughputOverride = f12elemf2
			}
			if f12iter.RegionName != nil {
				f12elem.RegionName = f12iter.RegionName
			}
			if f12iter.ReplicaStatus != nil {
				f12elem.ReplicaStatus = f12iter.ReplicaStatus
			}
			if f12iter.ReplicaStatusDescription != nil {
				f12elem.ReplicaStatusDescription = f12iter.ReplicaStatusDescription
			}
			if f12iter.ReplicaStatusPercentProgress != nil {
				f12elem.ReplicaStatusPercentProgress = f12iter.ReplicaStatusPercentProgress
			}
			f12 = append(f12, f12elem)
		}
		cr.Status.AtProvider.Replicas = f12
	}
	if resp.TableDescription.RestoreSummary != nil {
		f13 := &svcapitypes.RestoreSummary{}
		if resp.TableDescription.RestoreSummary.RestoreDateTime != nil {
			f13.RestoreDateTime = &metav1.Time{*resp.TableDescription.RestoreSummary.RestoreDateTime}
		}
		if resp.TableDescription.RestoreSummary.RestoreInProgress != nil {
			f13.RestoreInProgress = resp.TableDescription.RestoreSummary.RestoreInProgress
		}
		if resp.TableDescription.RestoreSummary.SourceBackupArn != nil {
			f13.SourceBackupARN = resp.TableDescription.RestoreSummary.SourceBackupArn
		}
		if resp.TableDescription.RestoreSummary.SourceTableArn != nil {
			f13.SourceTableARN = resp.TableDescription.RestoreSummary.SourceTableArn
		}
		cr.Status.AtProvider.RestoreSummary = f13
	}
	if resp.TableDescription.SSEDescription != nil {
		f14 := &svcapitypes.SSEDescription{}
		if resp.TableDescription.SSEDescription.InaccessibleEncryptionDateTime != nil {
			f14.InaccessibleEncryptionDateTime = &metav1.Time{*resp.TableDescription.SSEDescription.InaccessibleEncryptionDateTime}
		}
		if resp.TableDescription.SSEDescription.KMSMasterKeyArn != nil {
			f14.KMSMasterKeyARN = resp.TableDescription.SSEDescription.KMSMasterKeyArn
		}
		if resp.TableDescription.SSEDescription.SSEType != nil {
			f14.SSEType = resp.TableDescription.SSEDescription.SSEType
		}
		if resp.TableDescription.SSEDescription.Status != nil {
			f14.Status = resp.TableDescription.SSEDescription.Status
		}
		cr.Status.AtProvider.SSEDescription = f14
	}
	if resp.TableDescription.TableArn != nil {
		cr.Status.AtProvider.TableARN = resp.TableDescription.TableArn
	}
	if resp.TableDescription.TableId != nil {
		cr.Status.AtProvider.TableID = resp.TableDescription.TableId
	}
	if resp.TableDescription.TableName != nil {
		cr.Status.AtProvider.TableName = resp.TableDescription.TableName
	}
	if resp.TableDescription.TableSizeBytes != nil {
		cr.Status.AtProvider.TableSizeBytes = resp.TableDescription.TableSizeBytes
	}
	if resp.TableDescription.TableStatus != nil {
		cr.Status.AtProvider.TableStatus = resp.TableDescription.TableStatus
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*svcapitypes.Table)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}
	if err := e.preUpdate(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "pre-update failed")
	}
	return e.postUpdate(ctx, cr, managed.ExternalUpdate{}, nil)
}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.Table)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteTableInput(cr)
	_, err := e.client.DeleteTableWithContext(ctx, input)
	return errors.Wrap(cpresource.Ignore(IsNotFound, err), errDelete)
}
