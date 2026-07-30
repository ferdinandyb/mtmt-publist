# MTMT PubList

A small server that fetches and compiles relevant information from the [MTMT](https://mtmt.hu) MYCITE2 API in an easy-to-display format.

Previous queries are cached so a) MTMT is not bombarded for every page request and b) it has fallbacks for when MTMT is down. Journal titles are also cached globally (those should not really change, but querying the API for each of them takes a long time).

## Usage without hosting

Use the `etoservice-standalone.html` (if viewing on github, you can find it above in the file list view) as an example to use the service which is hosted at `etoservice.elte.hu`. Just copy the contents into page on your website.

An example how this would look is [here](https://bence.ferdinandy.com/publications/), bar styling (e.g. fonts) of course.

If `etoservice.elte.hu` is unreachable, the page falls back to querying the MTMT API directly in the browser. This works for any list, but for institutions with many publications it can take up to a minute and download several megabytes, since MTMT does not offer a lighter response format — the hosted service exists precisely to avoid that cost on every page view.

Depending on whether you want to generate a list for yourself (single author) or an institution, you'd need to make the below modifications to the file, at the top of the `<script>` block:

### Single author

```
const MODE = "user";
const MTIDS = ["10028021"];
```

`MTIDS` should hold the specific author's mtid. I usually get this from the URL of the author's page, e.g. if I search for myself I land on this URL: https://m2.mtmt.hu/gui2/?type=authors&mode=browse&sel=10028021, which has my mtmt id at the end.

### Institutions

Rather similar, but set `MODE` to `"institute"` and list one or more IDs, e.g. if your department also hosts some research groups which you want to handle together.

```
const MODE = "institute";
const MTIDS = ["338", "12724", "11351", "20298"];
```
