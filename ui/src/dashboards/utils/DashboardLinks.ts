import {getDashboards} from 'src/dashboards/apis'

import {AxiosResponse} from 'axios'
import {DashboardsResponse} from 'src/types/apis/dashboards'
import {Source} from 'src/types/sources'
import {Dashboard, DashboardSwitcherLink} from 'src/types/dashboards'

type DashboardsRequest = () => Promise<AxiosResponse<DashboardsResponse>>

class DashboardsLinks {
  public static DASHBOARDS_AJAX: DashboardsRequest = getDashboards as DashboardsRequest
  public static EMPTY: DashboardsLinks = new DashboardsLinks([], null)

  public static async load(
    dashboardsAJAX: DashboardsRequest,
    source: Source
  ): Promise<DashboardsLinks> {
    const {
      data: {dashboards},
    } = await dashboardsAJAX()

    return DashboardsLinks.fromDashboards(dashboards, source)
  }

  public static fromDashboards(
    dashboards: Dashboard[],
    source: Source
  ): DashboardsLinks {
    const links = dashboards.map(d => {
      return {
        key: String(d.id),
        text: d.name,
        to: `/sources/${source.id}/dashboards/${d.id}`,
      }
    })

    return new DashboardsLinks(links, null)
  }

  private links: DashboardSwitcherLink[]
  private active: DashboardSwitcherLink

  constructor(links: DashboardSwitcherLink[], active: DashboardSwitcherLink) {
    this.links = links
    this.active = active
  }

  public withActiveDashboard(dashboard: Dashboard) {
    if (!dashboard) {
      return new DashboardsLinks(this.links, null)
    }

    const active = this.links.find(link => link.key === String(dashboard.id))

    return new DashboardsLinks(this.links, active)
  }

  public toArray(): DashboardSwitcherLink[] {
    return this.links.map(link => ({...link, isActive: link === this.active}))
  }

  public get count() {
    return this.links.length
  }
}

export default DashboardsLinks
