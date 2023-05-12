package status

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testJsonResponse = `{
  "page": {
    "id": "kctbh9vrtdwd",
    "name": "GitHub",
    "url": "https://www.githubstatus.com",
    "time_zone": "Etc/UTC",
    "updated_at": "2023-05-12T07:59:13.397Z"
  },
  "components": [
    {
      "id": "8l4ygp009s5s",
      "name": "Git Operations",
      "status": "operational",
      "created_at": "2017-01-31T20:05:05.370Z",
      "updated_at": "2023-05-11T14:40:16.954Z",
      "position": 1,
      "description": "Performance of git clones, pulls, pushes, and associated operations",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "brv1bkgrwx7q",
      "name": "API Requests",
      "status": "operational",
      "created_at": "2017-01-31T20:01:46.621Z",
      "updated_at": "2023-05-11T14:40:15.561Z",
      "position": 2,
      "description": "Requests for GitHub APIs",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "4230lsnqdsld",
      "name": "Webhooks",
      "status": "operational",
      "created_at": "2019-11-13T18:00:24.256Z",
      "updated_at": "2023-05-11T14:40:18.323Z",
      "position": 3,
      "description": "Real time HTTP callbacks of user-generated and system events",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "0l2p9nhqnxpd",
      "name": "Visit www.githubstatus.com for more information",
      "status": "operational",
      "created_at": "2018-12-05T19:39:40.838Z",
      "updated_at": "2022-09-07T00:08:33.519Z",
      "position": 4,
      "description": null,
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "kr09ddfgbfsf",
      "name": "Issues",
      "status": "operational",
      "created_at": "2017-01-31T20:01:46.638Z",
      "updated_at": "2023-05-11T14:40:17.642Z",
      "position": 5,
      "description": "Requests for Issues on GitHub.com",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "hhtssxt0f5v2",
      "name": "Pull Requests",
      "status": "operational",
      "created_at": "2020-09-02T15:39:06.329Z",
      "updated_at": "2023-05-11T19:00:39.146Z",
      "position": 6,
      "description": "Requests for Pull Requests on GitHub.com",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "br0l2tvcx85d",
      "name": "Actions",
      "status": "operational",
      "created_at": "2019-11-13T18:02:19.432Z",
      "updated_at": "2023-05-11T14:40:14.904Z",
      "position": 7,
      "description": "Workflows, Compute and Orchestration for GitHub Actions",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "st3j38cctv9l",
      "name": "Packages",
      "status": "operational",
      "created_at": "2019-11-13T18:02:40.064Z",
      "updated_at": "2023-04-27T09:56:19.514Z",
      "position": 8,
      "description": "API requests and webhook delivery for GitHub Packages",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "vg70hn9s2tyj",
      "name": "Pages",
      "status": "operational",
      "created_at": "2017-01-31T20:04:33.923Z",
      "updated_at": "2023-05-11T14:46:14.253Z",
      "position": 9,
      "description": "Frontend application and API servers for Pages builds",
      "showcase": false,
      "start_date": null,
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "h2ftsgbw7kmk",
      "name": "Codespaces",
      "status": "operational",
      "created_at": "2021-08-11T16:02:09.505Z",
      "updated_at": "2023-05-11T14:40:16.240Z",
      "position": 10,
      "description": "Orchestration and Compute for GitHub Codespaces",
      "showcase": false,
      "start_date": "2021-08-11",
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    },
    {
      "id": "pjmpxvq2cmr2",
      "name": "Copilot",
      "status": "operational",
      "created_at": "2022-06-21T16:04:33.017Z",
      "updated_at": "2023-05-04T16:18:39.969Z",
      "position": 11,
      "description": null,
      "showcase": false,
      "start_date": "2022-06-21",
      "group_id": null,
      "page_id": "kctbh9vrtdwd",
      "group": false,
      "only_show_if_degraded": false
    }
  ]
}`
)

func TestVerifyConstants(t *testing.T) {
	require.Equal(t, "operational", COMPONENT_OPERATIONAL)
	require.Equal(t, "degraded_performance", COMPONENT_DEGREDADED_PERFORMANCE)
	require.Equal(t, "partial_outage", COMPONENT_PARTIAL_OUTAGE)
	require.Equal(t, "major_outage", COMPONENT_MAJOR_OUTAGE)
}
