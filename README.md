# MTMT PubList

A small server that fetches and compiles relevant information from the [MTMT](https://mtmt.hu) MYCITE2 API in an easy-to-display format.

## Usage without hosting

Use the `etoservice-standalone.html` as an example to use the service which is hosted at `etoservice.elte.hu`.

An example how this would look is [here](https://bence.ferdinandy.com/publications/), bar styling like fonts of course.

### Single author

To get the publications of one person change line 45

```
    fetch ("https://etoservice.elte.hu/mtmt-publist/user?mtid=10028021").then(res => res.json()).then(
```

so that the `mtid=XXXXXXX` has the specific author's mtid. I usually get this from the URL of the author's page, e.g. if I search for myself I land on this URL: https://m2.mtmt.hu/gui2/?type=authors&mode=browse&sel=10028021, which has my mtmt id at the end.

### Institutions

Rather similar, but you need to change `user` to `institute` in the url and here you can specify several ID-s, e.g. if your department also hosts some research groups which you want to handle together.

```
    fetch ("https://etoservice.elte.hu/mtmt-publist/institute?mtid=338&mtid=12724&mtid=11351&mtid=20298").then(res => res.json()).then(
```
