import React, {PureComponent} from 'react'
import {Link} from 'react-router'
import _ from 'lodash'
import classnames from 'classnames'

import OnClickOutside from 'src/shared/components/OnClickOutside'
import FancyScrollbar from 'src/shared/components/FancyScrollbar'

import {ErrorHandling} from 'src/shared/decorators/errors'

import {DROPDOWN_MENU_MAX_HEIGHT} from 'src/shared/constants/index'
import {getLinksWithActiveStatus} from 'src/dashboards/utils/DashboardLinks'
import {DashboardNameLinks, DashboardSwitcherLink} from 'src/types/dashboards'

interface Props {
  dashboardLinks: DashboardNameLinks
}

interface State {
  isOpen: boolean
}

@ErrorHandling
class DashboardSwitcher extends PureComponent<Props, State> {
  constructor(props) {
    super(props)

    this.state = {isOpen: false}
  }

  public render() {
    const {isOpen} = this.state

    const openClass = isOpen ? 'open' : ''

    return (
      <div className={`dropdown dashboard-switcher ${openClass}`}>
        <button
          className="btn btn-square btn-default btn-sm dropdown-toggle"
          onClick={this.handleToggleMenu}
        >
          <span className="icon dash-h" />
        </button>
        <ul className="dropdown-menu">
          <FancyScrollbar
            autoHeight={true}
            maxHeight={DROPDOWN_MENU_MAX_HEIGHT}
          >
            {this.links}
          </FancyScrollbar>
        </ul>
      </div>
    )
  }

  public handleClickOutside = () => {
    this.setState({isOpen: false})
  }

  private handleToggleMenu = () => {
    this.setState({isOpen: !this.state.isOpen})
  }

  private handleCloseMenu = () => {
    this.setState({isOpen: false})
  }

  private get links(): JSX.Element[] {
    return _.sortBy(this.dashboardNameLinks, ['text', 'key']).map(link => {
      return (
        <li
          key={link.key}
          className={classnames('dropdown-item', {active: link.isActive})}
        >
          <Link to={link.to} onClick={this.handleCloseMenu}>
            {link.text}
          </Link>
        </li>
      )
    })
  }

  private get dashboardNameLinks(): DashboardSwitcherLink[] {
    return getLinksWithActiveStatus(this.props.dashboardLinks)
  }
}

export default OnClickOutside(DashboardSwitcher)
