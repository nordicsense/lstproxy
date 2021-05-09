
data <- read.csv("/data/Research/Analysis/Fin-Meteo/Aqua-Day.csv")
stations <- read.csv("/data/Research/Data/Meteo/Finland/Fin_ListMeteoStations.csv")

x <- cbind(as.integer(names(split(data, data$id))), sapply(split(data, data$id), nrow))
stations <- cbind(stations, matches=x[match(stations$FIMIS_ID, x[,1]),2])
stations <- stations[,c(3,2,6,7,9)]

colnames(stations) <- c("id","name","lat","lon","matches")
write.csv(stations, file= "/data/Research/Data/Meteo/finmet/analysis/Finland/station-stats.csv", quote = FALSE, row.names=FALSE)
