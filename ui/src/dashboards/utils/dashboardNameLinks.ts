import {GetDashboards} from 'src/types/apis/dashboards'
import {Source} from 'src/types/sources'
import {Dashboard, DashboardNameLinks} from 'src/types/dashboards'

export const EMPTY_LINKS = {
  links: [],
  active: null,
}

export const loadDashboardLinks = async (
  dashboardsAJAX: GetDashboards,
  source: Source
): Promise<DashboardNameLinks> => {
  const {
    data: {dashboards},
  } = await dashboardsAJAX()

  return linksFromDashboards(dashboards, source)
}

const linksFromDashboards = (
  dashboards: Dashboard[],
  source: Source
): DashboardNameLinks => {
  const links = dashboards.map(d => {
    return {
      key: String(d.id),
      text: d.name,
      to: `/sources/${source.id}/dashboards/${d.id}`,
    }
  })

  return {links, active: null}
}

export const updateActiveDashboardLink = (
  dashboardLinks: DashboardNameLinks,
  dashboard: Dashboard
) => {
  if (!dashboard) {
    return {...dashboardLinks, active: null}
  }

  const active = dashboardLinks.links.find(
    link => link.key === String(dashboard.id)
  )

  return {...dashboardLinks, active}
}
