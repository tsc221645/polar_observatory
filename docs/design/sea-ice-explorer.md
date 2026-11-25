# Sea ice explorer wireframe

**Status** Proposed

```text
Polar Observatory                        Overview | Sea Ice | Data | Methodology

Southern Hemisphere sea ice extent
From [YYYY-MM-DD]   To [YYYY-MM-DD]   [Apply]

Latest observation       Source                     Data status
[Extent and UTC date]    [NOAA/NSIDC citation]       [Freshness and last successful sync]

Daily sea ice extent (km²)
[Time series with gaps for absent observations]
[View data table]
```

The first sea ice view uses URL date filters with an inclusive start and an exclusive end. The daily series shows gaps wherever dates or values are absent. The data table provides the same dates, units, quality flags, and source as the chart. The freshness status distinguishes recent observations from a recent successful sync.

On small screens, stack filters, status, and chart in reading order. Keep filter labels visible, provide keyboard access to the table, and convey each state in text. Cover loading, empty results, stale and partial data, and errors with an available retry. Display a prominent DEMO label whenever the view uses fixtures.
