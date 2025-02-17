package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CMSgov/dpc/attribution/attributiontest"
	"github.com/CMSgov/dpc/attribution/middleware"
	"github.com/CMSgov/dpc/attribution/model/v2"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GroupRepositoryTestSuite struct {
	suite.Suite
	fakeGrp *v2.Group
}

func (suite *GroupRepositoryTestSuite) SetupTest() {
	g := v2.Group{}
	err := faker.FakeData(&g)
	if err != nil {
		fmt.Printf("ERR %v\n", err)
	}
	var i v2.Info
	_ = json.Unmarshal([]byte(attributiontest.Groupjson), &i)
	g.Info = i
	suite.fakeGrp = &g
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}

func (suite *GroupRepositoryTestSuite) TestInsertErrorInRepo() {
	db, mock := newMock()
	defer db.Close()
	repo := NewGroupRepo(db)
	ctx := context.Background()

	expectedInsertQuery := `INSERT INTO group \(info\) VALUES \(\$1\) returning id, version, created_at, updated_at, info`

	rows := sqlmock.NewRows([]string{"id", "version", "created_at", "updated_at", "info"})

	mock.ExpectQuery(expectedInsertQuery).WithArgs(suite.fakeGrp.Info).WillReturnRows(rows)

	b, _ := json.Marshal(suite.fakeGrp.Info)
	group, err := repo.Insert(ctx, b)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), group)
}

func (suite *GroupRepositoryTestSuite) TestInsert() {
	db, mock := newMock()
	defer db.Close()
	repo := NewGroupRepo(db)
	ctx := context.WithValue(context.Background(), middleware.ContextKeyOrganization, "12345")

	expectedInsertQuery := `INSERT INTO "group" \(info, organization_id\) VALUES \(\$1, \$2\) returning id, version, created_at, updated_at, info, organization_id`

	rows := sqlmock.NewRows([]string{"id", "version", "created_at", "updated_at", "info", "organization_id"}).
		AddRow(suite.fakeGrp.ID, suite.fakeGrp.Version, suite.fakeGrp.CreatedAt, suite.fakeGrp.UpdatedAt, suite.fakeGrp.Info, suite.fakeGrp.OrganizationID)

	mock.ExpectQuery(expectedInsertQuery).WithArgs(suite.fakeGrp.Info, "12345").WillReturnRows(rows)

	b, _ := json.Marshal(suite.fakeGrp.Info)
	group, err := repo.Insert(ctx, b)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.fakeGrp.ID, group.ID)
}
