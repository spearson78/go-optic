package ohttp_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/exp/ohttp"
	"github.com/spearson78/go-optic/ojson"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestRestJSON(t *testing.T) {

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(testJson)),
			Header:     make(http.Header),
		}
	})

	if ret, ok, err := GetFirst(Compose(RestJSON[any, ReturnOne, Pure](client, 8192, nil), ojson.Nth(0).Key("commit").Key("author").Key("name")), "http://localhost:8080/test"); !ok || err != nil || ret != "dependabot[bot]" {
		t.Fatal(ret, ok, err)
	}

}

const testJson = `
[
  {
    "sha": "588ff1874c8c394253c231733047a550efe78260",
    "node_id": "C_kwDOAE3WVdoAKDU4OGZmMTg3NGM4YzM5NDI1M2MyMzE3MzMwNDdhNTUwZWZlNzgyNjA",
    "commit": {
      "author": {
        "name": "dependabot[bot]",
        "email": "49699333+dependabot[bot]@users.noreply.github.com",
        "date": "2024-12-29T12:54:10Z"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "date": "2024-12-29T12:54:10Z"
      },
      "message": "build(deps): bump jinja2 from 3.1.4 to 3.1.5 in /docs (#3226)\n\nBumps [jinja2](https://github.com/pallets/jinja) from 3.1.4 to 3.1.5.\r\n- [Release notes](https://github.com/pallets/jinja/releases)\r\n- [Changelog](https://github.com/pallets/jinja/blob/main/CHANGES.rst)\r\n- [Commits](https://github.com/pallets/jinja/compare/3.1.4...3.1.5)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: jinja2\r\n  dependency-type: direct:production\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
      "tree": {
        "sha": "e47dd35840a08fc1fa45e364ad536ac0941c1f72",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/e47dd35840a08fc1fa45e364ad536ac0941c1f72"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/588ff1874c8c394253c231733047a550efe78260",
      "comment_count": 0,
      "verification": {
        "verified": true,
        "reason": "valid",
        "signature": "-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJncUZyCRC1aQ7uu5UhlAAAkqsQAGDitC5DSe3ySbRtmHzPf7iv\n5Khp9i9SyBWrRkD1Craega8o9YJnmkg2tMl1ewZ0ivRs95OIhJk7RDbnNWxOupmJ\nfvLXbi2KHqz4IJb3gzHmZ0k2mlPI/iGlSbvIuOs4ot7hfZP/4gXZgIPoMMkNJUu1\n+hgb8RtiKpQujamxNytUaahW/vTIO5OR+iinfB0fDk66LDiLtNYWlZeOVoGiD5rD\nogha/Cm7InhdiqUiac/L2fezPbi1ZDZfn2iKUHhVyeHb5RVafU20vzHw7gk32mFz\n9PgQcnliHrG1+l9dLSNp+7b/Lft57VDE8p8C8IJEsjwlXXChAOJn0x5Nrl+mo8Tm\ngt/Q18DaRYyzwHUdSMX9Kp4bTnPoydMsriT7aibFGgkI/VLlK5/zAiUOSw0ffSmV\nrF/DufMxWqkUy2U5egGbyo0ZCFu+198jwcQCQC4p+09Z9xB4/YkR1jgFmoJB6ZGS\nswj8SXN6dATpI0GFyCPka9cIZU780Q6RKWd3MDBJqAfth5lRVEkGRcZB0k3wL6B4\n+O870LGWn/eqyDxniWvYEz9MKTZhuT46FGS8d7WWsOdMBy+ssO78q9/XHy96z25w\najkM4XErrpnGtxH1CtYt4SkYBB9oTdY2pgJDnXQcKiUKScTgQPoqQ4OfH0R27xmX\n1anb5H4HAqQj3CrD9qe2\n=Zja6\n-----END PGP SIGNATURE-----\n",
        "payload": "tree e47dd35840a08fc1fa45e364ad536ac0941c1f72\nparent bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195\nauthor dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com> 1735476850 +0900\ncommitter GitHub <noreply@github.com> 1735476850 +0900\n\nbuild(deps): bump jinja2 from 3.1.4 to 3.1.5 in /docs (#3226)\n\nBumps [jinja2](https://github.com/pallets/jinja) from 3.1.4 to 3.1.5.\r\n- [Release notes](https://github.com/pallets/jinja/releases)\r\n- [Changelog](https://github.com/pallets/jinja/blob/main/CHANGES.rst)\r\n- [Commits](https://github.com/pallets/jinja/compare/3.1.4...3.1.5)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: jinja2\r\n  dependency-type: direct:production\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
        "verified_at": "2024-12-29T12:54:13Z"
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/588ff1874c8c394253c231733047a550efe78260",
    "html_url": "https://github.com/jqlang/jq/commit/588ff1874c8c394253c231733047a550efe78260",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/588ff1874c8c394253c231733047a550efe78260/comments",
    "author": {
      "login": "dependabot[bot]",
      "id": 49699333,
      "node_id": "MDM6Qm90NDk2OTkzMzM=",
      "avatar_url": "https://avatars.githubusercontent.com/in/29110?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/dependabot%5Bbot%5D",
      "html_url": "https://github.com/apps/dependabot",
      "followers_url": "https://api.github.com/users/dependabot%5Bbot%5D/followers",
      "following_url": "https://api.github.com/users/dependabot%5Bbot%5D/following{/other_user}",
      "gists_url": "https://api.github.com/users/dependabot%5Bbot%5D/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/dependabot%5Bbot%5D/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/dependabot%5Bbot%5D/subscriptions",
      "organizations_url": "https://api.github.com/users/dependabot%5Bbot%5D/orgs",
      "repos_url": "https://api.github.com/users/dependabot%5Bbot%5D/repos",
      "events_url": "https://api.github.com/users/dependabot%5Bbot%5D/events{/privacy}",
      "received_events_url": "https://api.github.com/users/dependabot%5Bbot%5D/received_events",
      "type": "Bot",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "web-flow",
      "id": 19864447,
      "node_id": "MDQ6VXNlcjE5ODY0NDQ3",
      "avatar_url": "https://avatars.githubusercontent.com/u/19864447?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/web-flow",
      "html_url": "https://github.com/web-flow",
      "followers_url": "https://api.github.com/users/web-flow/followers",
      "following_url": "https://api.github.com/users/web-flow/following{/other_user}",
      "gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
      "organizations_url": "https://api.github.com/users/web-flow/orgs",
      "repos_url": "https://api.github.com/users/web-flow/repos",
      "events_url": "https://api.github.com/users/web-flow/events{/privacy}",
      "received_events_url": "https://api.github.com/users/web-flow/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
        "url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
        "html_url": "https://github.com/jqlang/jq/commit/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195"
      }
    ]
  },
  {
    "sha": "bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "node_id": "C_kwDOAE3WVdoAKGJjYmYyYjQ2MTY4OTNjZjJlNmE0ZjhhOTJkYzNkYjNiMWVlYjExOTU",
    "commit": {
      "author": {
        "name": "lectrical",
        "email": "14344693+lectrical@users.noreply.github.com",
        "date": "2024-12-29T12:53:16Z"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "date": "2024-12-29T12:53:16Z"
      },
      "message": "Generate provenance attestations for release artifacts and docker image (#3225)\n\nAdding https://github.com/actions/attest-build-provenance to the ci builds so\r\nthat the release assets and docker image for the next release tag generate\r\nsigned build provenance attestations for workflow artifacts.",
      "tree": {
        "sha": "3b91252273d4a46a78a74974a09cb3dd62d73223",
        "url": "https://api.github.com/repos/jqlang/jq/git/trees/3b91252273d4a46a78a74974a09cb3dd62d73223"
      },
      "url": "https://api.github.com/repos/jqlang/jq/git/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
      "comment_count": 0,
      "verification": {
        "verified": true,
        "reason": "valid",
        "signature": "-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJncUY8CRC1aQ7uu5UhlAAAzxAQABJceWBtoI3mfLmZlUa5L5s5\ne98L6+EPWucJlfTHsIGKYLZbK1gITyNVzzgLaizo9+ht5cm+2I1H+nqcbhIYg5ge\nM0w838W6EzkF8EKTMOElI6YQuVfTZEgIu4nlF9e6484kkBm4ed1d1iLk5NpLa+6a\nOjNXFM8gbBNj/4/+SZD5NJkNp/An53JrKA/NQkzw9EINmnOfkQ3Og2NKxm74fUGu\nbnJS2oA61yUZFq+4Po8hudOutNkCMpm+MpHs3t/E/lFXMdgoiTjSmnCEswll3Rcj\ngKYxfrLCtOkC3P9IsN5OuYfQs3vIZBaeKxwMOLbrnarG2W4kgMwkF/fM1ToaUkrb\nEsepBoW4SLe118bwC1g8syWJX8l8Pd0YXJlOil1Og828p9oAGhUP/v3ycUL3sYaZ\ntOFEcyEJgZJgEcLpMZ0QXySWIyp5GGLb/keMdD6HxmnjLM9cliTNYs4YbzJoveo1\n3lykmA08chqldFADpJM93pWUtLOyV2WjAF2IdDBZQ/ZCni9Ge6a8v3JXR89jPHV/\nJhbzkoSUCR40IH7+ReAiRiz2MYLUUQMS2CHrwftCZGvIcYz7wrpSXDxMHcVrSdr4\n0bjbVoDeCVzJ9IhJR705ZglFh1+KvwaquO3fT8r/6wFITHAaw9u9EmMIJN7IfMnt\ni0a+sTL6HIqy8KIwK3Pk\n=ByY+\n-----END PGP SIGNATURE-----\n",
        "payload": "tree 3b91252273d4a46a78a74974a09cb3dd62d73223\nparent 8bcdc9304ace5f2cc9bf662ab8998d75537e05f0\nauthor lectrical <14344693+lectrical@users.noreply.github.com> 1735476796 +0000\ncommitter GitHub <noreply@github.com> 1735476796 +0900\n\nGenerate provenance attestations for release artifacts and docker image (#3225)\n\nAdding https://github.com/actions/attest-build-provenance to the ci builds so\r\nthat the release assets and docker image for the next release tag generate\r\nsigned build provenance attestations for workflow artifacts.",
        "verified_at": "2024-12-29T12:53:19Z"
      }
    },
    "url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "html_url": "https://github.com/jqlang/jq/commit/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195",
    "comments_url": "https://api.github.com/repos/jqlang/jq/commits/bcbf2b4616893cf2e6a4f8a92dc3db3b1eeb1195/comments",
    "author": {
      "login": "lectrical",
      "id": 14344693,
      "node_id": "MDQ6VXNlcjE0MzQ0Njkz",
      "avatar_url": "https://avatars.githubusercontent.com/u/14344693?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/lectrical",
      "html_url": "https://github.com/lectrical",
      "followers_url": "https://api.github.com/users/lectrical/followers",
      "following_url": "https://api.github.com/users/lectrical/following{/other_user}",
      "gists_url": "https://api.github.com/users/lectrical/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/lectrical/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/lectrical/subscriptions",
      "organizations_url": "https://api.github.com/users/lectrical/orgs",
      "repos_url": "https://api.github.com/users/lectrical/repos",
      "events_url": "https://api.github.com/users/lectrical/events{/privacy}",
      "received_events_url": "https://api.github.com/users/lectrical/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "committer": {
      "login": "web-flow",
      "id": 19864447,
      "node_id": "MDQ6VXNlcjE5ODY0NDQ3",
      "avatar_url": "https://avatars.githubusercontent.com/u/19864447?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/web-flow",
      "html_url": "https://github.com/web-flow",
      "followers_url": "https://api.github.com/users/web-flow/followers",
      "following_url": "https://api.github.com/users/web-flow/following{/other_user}",
      "gists_url": "https://api.github.com/users/web-flow/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/web-flow/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/web-flow/subscriptions",
      "organizations_url": "https://api.github.com/users/web-flow/orgs",
      "repos_url": "https://api.github.com/users/web-flow/repos",
      "events_url": "https://api.github.com/users/web-flow/events{/privacy}",
      "received_events_url": "https://api.github.com/users/web-flow/received_events",
      "type": "User",
      "user_view_type": "public",
      "site_admin": false
    },
    "parents": [
      {
        "sha": "8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
        "url": "https://api.github.com/repos/jqlang/jq/commits/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0",
        "html_url": "https://github.com/jqlang/jq/commit/8bcdc9304ace5f2cc9bf662ab8998d75537e05f0"
      }
    ]
  }
]	
`
