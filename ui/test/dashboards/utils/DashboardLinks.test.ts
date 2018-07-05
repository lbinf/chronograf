import DashboardLinks from 'src/dashboards/utils/dashboardLinks'
import {dashboard, source} from 'test/resources'

describe('dashboards.utils.DashboardLinks', () => {
  describe('.load', () => {
    const socure = {...source, id: '897'}

    const dashboards = [
      {
        ...dashboard,
        id: 123,
        name: 'Test Dashboard',
      },
    ]

    const data = {
      dashboards,
    }

    const axiosResponse = {
      data,
      status: 200,
      statusText: 'Okay',
      headers: null,
      config: null,
    }

    const getDashboards = async () => axiosResponse

    it('can load dashboard links for source', async () => {
      const actualLinks = await DashboardLinks.load(getDashboards, socure)

      const expectedLinks = new DashboardLinks(
        [
          {
            key: '123',
            text: 'Test Dashboard',
            to: '/sources/897/dashboards/123',
          },
        ],
        null
      )

      expect(actualLinks).toEqual(expectedLinks)
      expect(actualLinks.count).toBe(1)
    })
  })

  const activeDashboard = {
    ...dashboard,
    id: 123,
    name: 'Test Dashboard',
  }

  const activeLink = {
    key: '123',
    text: 'Test Dashboard',
    to: '/sources/897/dashboards/123',
  }

  const link1 = {
    key: '9001',
    text: 'Low Dash',
    to: '/sources/897/dashboards/9001',
  }

  const link2 = {
    key: '2282',
    text: 'Low Dash',
    to: '/sources/897/dashboards/2282',
  }

  const links = [link1, activeLink, link2]

  describe('#withActiveDashboard', () => {
    it('can set the active link', () => {
      const loadedLinks = new DashboardLinks(links, null)
      const actualLinks = loadedLinks.withActiveDashboard(activeDashboard)
      const expectedLinks = new DashboardLinks(links, activeLink)

      expect(actualLinks).toEqual(expectedLinks)
    })

    it('can handle a missing dashboard', () => {
      const loadedLinks = new DashboardLinks(links, activeLink)
      const actualLinks = loadedLinks.withActiveDashboard(undefined)
      const expectedLinks = new DashboardLinks(links, null)

      expect(actualLinks).toEqual(expectedLinks)
    })
  })

  it('can convert to an array with active status', () => {
    const loadedLinks = new DashboardLinks(links, activeLink)
    const actualArray = loadedLinks.toArray()

    const expectedArray = [
      {
        ...link1,
        isActive: false,
      },
      {
        ...activeLink,
        isActive: true,
      },
      {
        ...link2,
        isActive: false,
      },
    ]

    expect(actualArray).toEqual(expectedArray)
  })
})
