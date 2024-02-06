# Docs

## Endpoints Used

### JSON Endpoints
| Judge | Endpoint |
| :- | :- |
| Codechef | [https://www.codechef.com/api/list/contests/all](https://www.codechef.com/api/list/contests/all) |
| Codeforces | [https://codeforces.com/api/contest.list](https://codeforces.com/api/contest.list) |
| Topcoder | [https://api.topcoder.com/v5/challenges/?status=Active&isLightweight=true&perPage=100&tracks%5B%5D=Dev&tracks%5B%5D=Des&tracks%5B%5D=DS&tracks%5B%5D=QA&types%5B%5D=CH&types%5B%5D=F2F&types%5B%5D=TSK](https://api.topcoder.com/v5/challenges/?status=Active&isLightweight=true&perPage=100&tracks%5B%5D=Dev&tracks%5B%5D=Des&tracks%5B%5D=DS&tracks%5B%5D=QA&types%5B%5D=CH&types%5B%5D=F2F&types%5B%5D=TSK) |
| Hackerank | [https://www.hackerrank.com/rest/contests/upcoming?limit=100](https://www.hackerrank.com/rest/contests/upcoming?limit=100) |
| Hackerearth | [https://www.hackerearth.com/chrome-extension/events/](https://www.hackerearth.com/chrome-extension/events/) |
| LeetCode | [https://leetcode.com/graphql?query=%7B%20allContests%20%7B%20title%20titleSlug%20startTime%20duration%20__typename%20%7D%20%7D](https://leetcode.com/graphql?query=%7B%20allContests%20%7B%20title%20titleSlug%20startTime%20duration%20__typename%20%7D%20%7D) |
| CsAcademy | [https://csacademy.com/contests](https://csacademy.com/contests) |

**NOTE**
Set `X-Requested-With: XMLHttpRequest` header on csacademy link to get json response

### Non-JSON Endpoints
| Judge | Endpoint |
| :- | :- |
| AtCoder | [https://atcoder.jp/contests](https://atcoder.jp/contests) |
| Toph | [https://toph.co/contests](https://toph.co/contests) |
