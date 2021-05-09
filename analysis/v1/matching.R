source("functions.R")

root <- "/data/Research/Data/Meteo/finmet/analysis/Finland"
datasets <- c("Aqua-Day", "Aqua-Night", "Terra-Day", "Terra-Night")

data <- readall(datasets)

pdf(sprintf("%s/ground-vs-lst.pdf", root), paper="a4", width=8, height=11.3)

# plotTitlePage("Observation times", "distributions")
# plotTimeDistributions(data)

# plotTitlePage("LST ~ T", "original LST values\nafter removing bad data")
plot2x2Density(data, colNames=c("LST", "t"))

dev.off()


#plotTitlePage("LST ~ T (corrected)", "All data:\nseasonally corrected\n(day/night separately)")
#plot2x2Density(data, colNames=COL_NAMES)

# plotTitlePage("LST ~ T", "by season\noriginal LST values\nafter removing bad data")
# plot2x2Density(splitBySeason(data)[["Aqua-Day"]], colNames=c("lst", "t"))
# plot2x2Density(splitBySeason(data)[["Terra-Day"]], colNames=c("lst", "t"))
# plot2x2Density(splitBySeason(data)[["Aqua-Night"]], colNames=c("lst", "t"))
# plot2x2Density(splitBySeason(data)[["Terra-Night"]], colNames=c("lst", "t"))

#plotTitlePage("LST ~ T (corrected)", "By season:\nseasonally corrected\n(day/night separately)")
#plot2x2Density(splitBySeason(data)[["Aqua-Day"]], colNames=COL_NAMES)
#plot2x2Density(splitBySeason(data)[["Terra-Day"]], colNames=COL_NAMES)
#plot2x2Density(splitBySeason(data)[["Aqua-Night"]], colNames=COL_NAMES)
#plot2x2Density(splitBySeason(data)[["Terra-Night"]], colNames=COL_NAMES)

#plotTitlePage("LST ~ T", "By month")
#plot3x4Density(splitByMonth(data)[["Aqua-Day"]][c(12,1:11)], colNames=COL_NAMES)
#plot3x4Density(splitByMonth(data)[["Terra-Day"]][c(12,1:11)], colNames=COL_NAMES)
#plot3x4Density(splitByMonth(data)[["Aqua-Night"]][c(12,1:11)], colNames=COL_NAMES)
#plot3x4Density(splitByMonth(data)[["Terra-Night"]][c(12,1:11)], colNames=COL_NAMES)

#plotTitlePage("LST ~ T", "By station (20 stations)")
#plot4x5Density(splitByStation(data)[["Aqua-Day"]][-1], colNames=COL_NAMES) # plot can only handle 20, there are 21
#plot4x5Density(splitByStation(data)[["Terra-Day"]][-1], colNames=COL_NAMES)
#plot4x5Density(splitByStation(data)[["Aqua-Night"]][-1], colNames=COL_NAMES)
#plot4x5Density(splitByStation(data)[["Terra-Night"]][-1], colNames=COL_NAMES)


# dev.off()

