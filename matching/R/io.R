read.matching <- function(fname) {
  res <- read.csv(fname, header=TRUE, stringsAsFactors=FALSE)
  res$datetime <- as.POSIXct(res$datetime, format="%Y-%m-%dT%H:%M:%OS", tz="UTC")
  res
}
