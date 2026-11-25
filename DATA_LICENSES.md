# Data use and attribution

Source review started on 2026-11-25. All entries in the [source catalog](config/sources.yaml) are disabled pending product review. Keep provider files outside Git. Fixtures committed during development must be explicitly identified as synthetic or have documented redistribution approval.

## NOAA/NSIDC Sea Ice Index G02135 v4

NSIDC requires a dataset citation. Include the dataset version, the subset used, and the actual access date when publishing derived observations. Use the provider's citation page for the full author list and citation format. [Dataset citation and acknowledgment](https://nsidc.org/data/g02135/versions/4)

The initial subset is the Southern Hemisphere daily extent CSV. A bounded read on 2026-11-25 verified the resource, its column header, a second row declaring units, and dated records beginning on 1978-10-26. Extent and Missing are both expressed in million square kilometers. The early dates in the sample are two days apart, so polling cadence cannot be used to infer missing historical rows. Review the provider terms for raw and derived redistribution, record any access limits, and check how the full series encodes missing data before enabling ingestion. [Data and image archive](https://nsidc.org/data/seaice_index/data-and-image-archive) · [Version 4 user guide](https://nsidc.org/data/documentation/sea-ice-index-version-4-user-guide)

## Later sources

Review the conditions, citation requirements, and cadence per selected AMRDC station product and per NOAA OISST v2.1 product before implementing either adapter. [AMRDC AWS collection](https://amrdcdata.ssec.wisc.edu/group/about/automatic-weather-station-project) · [NOAA OISST product](https://www.ncei.noaa.gov/products/optimum-interpolation-sst)
