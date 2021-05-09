
get.stations <- function() {
  fname <- system.file("data", "stations.csv", package="matching")
  read.csv(fname, stringsAsFactors = FALSE)
}


