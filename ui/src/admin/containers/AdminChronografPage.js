import React, {Component} from 'react'

import SourceIndicator from 'shared/components/SourceIndicator'
import AllUsersTable from 'src/admin/components/chronograf/AllUsersTable'

import FancyScrollbar from 'shared/components/FancyScrollbar'

import {DUMMY_USERS} from 'src/admin/constants/dummyUsers'

class AdminChronografPage extends Component {
  constructor(props) {
    super(props)

    this.state = {
      organizationName: null,
      users: DUMMY_USERS,
      selectedUsers: [],
    }
  }

  isSameUser = (userA, userB) => {
    return userA.name === userB.name && userA.provider === userB.provider
  }

  handleViewOrg = organizationName => () => {
    this.handleDeselectAllUsers()
    this.setState({organizationName})
  }

  handleToggleUserSelected = user => e => {
    e.preventDefault()

    const {selectedUsers} = this.state

    const isUserSelected = selectedUsers.find(u => this.isSameUser(user, u))

    const newSelectedUsers = isUserSelected
      ? selectedUsers.filter(u => !this.isSameUser(user, u))
      : [...selectedUsers, user]

    this.setState({selectedUsers: newSelectedUsers})
  }

  handleSelectAllUsers = () => {
    this.state.users.forEach(user => {
      user.selected = true
    })
  }
  handleDeselectAllUsers = () => {
    this.state.users.forEach(user => {
      user.selected = false
    })
  }

  render() {
    const {organizationName, users, selectedUsers} = this.state

    return (
      <div className="page">
        <div className="page-header">
          <div className="page-header__container">
            <div className="page-header__left">
              <h1 className="page-header__title">Chronograf Admin</h1>
            </div>
            <div className="page-header__right">
              <SourceIndicator />
              <button className="btn btn-primary btn-sm">
                <span className="icon plus" />
                Create Organization
              </button>
            </div>
          </div>
        </div>
        <FancyScrollbar className="page-contents">
          {users
            ? <div className="container-fluid">
                <div className="row">
                  <div className="col-xs-12">
                    <div className="panel panel-minimal">
                      <div className="panel-heading">
                        <h2 className="panel-title">
                          {organizationName
                            ? `Users in ${organizationName}`
                            : 'Users'}
                        </h2>
                      </div>
                      <div className="panel-body">
                        <AllUsersTable
                          users={users}
                          organizationName={organizationName}
                          onViewOrg={this.handleViewOrg}
                          onToggleUserSelected={this.handleToggleUserSelected}
                          selectedUsers={selectedUsers}
                          isSameUser={this.isSameUser}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            : <div className="page-spinner" />}
        </FancyScrollbar>
      </div>
    )
  }
}

// const {} = PropTypes

// AdminChronografPage.propTypes = {}
export default AdminChronografPage