// Copyright 2016 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// Package faketreeeval provides fake implementations of tree eval interfaces.
package faketreeeval

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgnotice"
	"github.com/cockroachdb/cockroach/pkg/sql/roleoption"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq/oid"
)

// DummySequenceOperators implements the tree.SequenceOperators interface by
// returning errors.
type DummySequenceOperators struct{}

var _ tree.EvalDatabase = &DummySequenceOperators{}

var errSequenceOperators = unimplemented.NewWithIssue(42508,
	"cannot evaluate scalar expressions containing sequence operations in this context")

// GetSerialSequenceNameFromColumn is part of the tree.SequenceOperators interface.
func (so *DummySequenceOperators) GetSerialSequenceNameFromColumn(
	ctx context.Context, tn *tree.TableName, columnName tree.Name,
) (*tree.TableName, error) {
	return nil, errors.WithStack(errSequenceOperators)
}

// ParseQualifiedTableName is part of the tree.EvalDatabase interface.
func (so *DummySequenceOperators) ParseQualifiedTableName(sql string) (*tree.TableName, error) {
	return nil, errors.WithStack(errSequenceOperators)
}

// ResolveTableName is part of the tree.EvalDatabase interface.
func (so *DummySequenceOperators) ResolveTableName(
	ctx context.Context, tn *tree.TableName,
) (tree.ID, error) {
	return 0, errors.WithStack(errSequenceOperators)
}

// LookupSchema is part of the tree.EvalDatabase interface.
func (so *DummySequenceOperators) LookupSchema(
	ctx context.Context, dbName, scName string,
) (bool, tree.SchemaMeta, error) {
	return false, nil, errors.WithStack(errSequenceOperators)
}

// IncrementSequence is part of the tree.SequenceOperators interface.
func (so *DummySequenceOperators) IncrementSequence(
	ctx context.Context, seqName *tree.TableName,
) (int64, error) {
	return 0, errors.WithStack(errSequenceOperators)
}

// GetLatestValueInSessionForSequence implements the tree.SequenceOperators
// interface.
func (so *DummySequenceOperators) GetLatestValueInSessionForSequence(
	ctx context.Context, seqName *tree.TableName,
) (int64, error) {
	return 0, errors.WithStack(errSequenceOperators)
}

// SetSequenceValue implements the tree.SequenceOperators interface.
func (so *DummySequenceOperators) SetSequenceValue(
	ctx context.Context, seqName *tree.TableName, newVal int64, isCalled bool,
) error {
	return errors.WithStack(errSequenceOperators)
}

// DummyEvalPlanner implements the tree.EvalPlanner interface by returning
// errors.
type DummyEvalPlanner struct{}

// UnsafeUpsertDescriptor is part of the EvalPlanner interface.
func (ep *DummyEvalPlanner) UnsafeUpsertDescriptor(
	ctx context.Context, descID int64, encodedDescriptor []byte, force bool,
) error {
	return errors.WithStack(errEvalPlanner)
}

// UnsafeDeleteDescriptor is part of the EvalPlanner interface.
func (ep *DummyEvalPlanner) UnsafeDeleteDescriptor(
	ctx context.Context, descID int64, force bool,
) error {
	return errors.WithStack(errEvalPlanner)
}

// UnsafeUpsertNamespaceEntry is part of the EvalPlanner interface.
func (ep *DummyEvalPlanner) UnsafeUpsertNamespaceEntry(
	ctx context.Context, parentID, parentSchemaID int64, name string, descID int64, force bool,
) error {
	return errors.WithStack(errEvalPlanner)
}

// UnsafeDeleteNamespaceEntry is part of the EvalPlanner interface.
func (ep *DummyEvalPlanner) UnsafeDeleteNamespaceEntry(
	ctx context.Context, parentID, parentSchemaID int64, name string, descID int64, force bool,
) error {
	return errors.WithStack(errEvalPlanner)
}

// CompactEngineSpan is part of the EvalPlanner interface.
func (ep *DummyEvalPlanner) CompactEngineSpan(
	ctx context.Context, nodeID int32, storeID int32, startKey []byte, endKey []byte,
) error {
	return errors.WithStack(errEvalPlanner)
}

var _ tree.EvalPlanner = &DummyEvalPlanner{}

var errEvalPlanner = pgerror.New(pgcode.ScalarOperationCannotRunWithoutFullSessionContext,
	"cannot evaluate scalar expressions using table lookups in this context")

// ParseQualifiedTableName is part of the tree.EvalDatabase interface.
func (ep *DummyEvalPlanner) ParseQualifiedTableName(sql string) (*tree.TableName, error) {
	return parser.ParseQualifiedTableName(sql)
}

// LookupSchema is part of the tree.EvalDatabase interface.
func (ep *DummyEvalPlanner) LookupSchema(
	ctx context.Context, dbName, scName string,
) (bool, tree.SchemaMeta, error) {
	return false, nil, errors.WithStack(errEvalPlanner)
}

// ResolveTableName is part of the tree.EvalDatabase interface.
func (ep *DummyEvalPlanner) ResolveTableName(
	ctx context.Context, tn *tree.TableName,
) (tree.ID, error) {
	return 0, errors.WithStack(errEvalPlanner)
}

// GetTypeFromValidSQLSyntax is part of the tree.EvalPlanner interface.
func (ep *DummyEvalPlanner) GetTypeFromValidSQLSyntax(sql string) (*types.T, error) {
	return nil, errors.WithStack(errEvalPlanner)
}

// EvalSubquery is part of the tree.EvalPlanner interface.
func (ep *DummyEvalPlanner) EvalSubquery(expr *tree.Subquery) (tree.Datum, error) {
	return nil, errors.WithStack(errEvalPlanner)
}

// ResolveTypeByOID implements the tree.TypeReferenceResolver interface.
func (ep *DummyEvalPlanner) ResolveTypeByOID(_ context.Context, _ oid.Oid) (*types.T, error) {
	return nil, errors.WithStack(errEvalPlanner)
}

// ResolveType implements the tree.TypeReferenceResolver interface.
func (ep *DummyEvalPlanner) ResolveType(
	_ context.Context, _ *tree.UnresolvedObjectName,
) (*types.T, error) {
	return nil, errors.WithStack(errEvalPlanner)
}

// DummyPrivilegedAccessor implements the tree.PrivilegedAccessor interface by returning errors.
type DummyPrivilegedAccessor struct{}

var _ tree.PrivilegedAccessor = &DummyPrivilegedAccessor{}

var errEvalPrivileged = pgerror.New(pgcode.ScalarOperationCannotRunWithoutFullSessionContext,
	"cannot evaluate privileged expressions in this context")

// LookupNamespaceID is part of the tree.PrivilegedAccessor interface.
func (ep *DummyPrivilegedAccessor) LookupNamespaceID(
	ctx context.Context, parentID int64, name string,
) (tree.DInt, bool, error) {
	return 0, false, errors.WithStack(errEvalPrivileged)
}

// LookupZoneConfigByNamespaceID is part of the tree.PrivilegedAccessor interface.
func (ep *DummyPrivilegedAccessor) LookupZoneConfigByNamespaceID(
	ctx context.Context, id int64,
) (tree.DBytes, bool, error) {
	return "", false, errors.WithStack(errEvalPrivileged)
}

// DummySessionAccessor implements the tree.EvalSessionAccessor interface by returning errors.
type DummySessionAccessor struct{}

var _ tree.EvalSessionAccessor = &DummySessionAccessor{}

var errEvalSessionVar = pgerror.New(pgcode.ScalarOperationCannotRunWithoutFullSessionContext,
	"cannot evaluate scalar expressions that access session variables in this context")

// GetSessionVar is part of the tree.EvalSessionAccessor interface.
func (ep *DummySessionAccessor) GetSessionVar(
	_ context.Context, _ string, _ bool,
) (bool, string, error) {
	return false, "", errors.WithStack(errEvalSessionVar)
}

// SetSessionVar is part of the tree.EvalSessionAccessor interface.
func (ep *DummySessionAccessor) SetSessionVar(_ context.Context, _, _ string) error {
	return errors.WithStack(errEvalSessionVar)
}

// HasAdminRole is part of the tree.EvalSessionAccessor interface.
func (ep *DummySessionAccessor) HasAdminRole(_ context.Context) (bool, error) {
	return false, errors.WithStack(errEvalSessionVar)
}

// HasRoleOption is part of the tree.EvalSessionAccessor interface.
func (ep *DummySessionAccessor) HasRoleOption(
	ctx context.Context, roleOption roleoption.Option,
) (bool, error) {
	return false, errors.WithStack(errEvalSessionVar)
}

// DummyClientNoticeSender implements the tree.ClientNoticeSender interface.
type DummyClientNoticeSender struct{}

var _ tree.ClientNoticeSender = &DummyClientNoticeSender{}

// BufferClientNotice is part of the tree.ClientNoticeSender interface.
func (c *DummyClientNoticeSender) BufferClientNotice(context.Context, pgnotice.Notice) {}

// DummyTenantOperator implements the tree.TenantOperator interface.
type DummyTenantOperator struct{}

var _ tree.TenantOperator = &DummyTenantOperator{}

var errEvalTenant = pgerror.New(pgcode.ScalarOperationCannotRunWithoutFullSessionContext,
	"cannot evaluate tenant operation in this context")

// CreateTenant is part of the tree.TenantOperator interface.
func (c *DummyTenantOperator) CreateTenant(_ context.Context, _ uint64) error {
	return errors.WithStack(errEvalTenant)
}

// DestroyTenant is part of the tree.TenantOperator interface.
func (c *DummyTenantOperator) DestroyTenant(_ context.Context, _ uint64) error {
	return errors.WithStack(errEvalTenant)
}

// GCTenant is part of the tree.TenantOperator interface.
func (c *DummyTenantOperator) GCTenant(_ context.Context, _ uint64) error {
	return errors.WithStack(errEvalTenant)
}
