
readone <- function(ds) {
  fname <- paste0(root, "/", ds, ".csv")
  x <- read.csv(fname, header=TRUE, stringsAsFactors=FALSE)
  x$datetime <- as.POSIXct(x$datetime,format="%Y-%m-%dT%H:%M:%OS", tz="UTC")
  x
}

readall <- function(datasets) {
  x <- lapply(datasets, readone)
  names(x) <- datasets
  x
}

newPage <- function() {
  dev.next()
  par(mfrow=c(1,1))
}

splitByStation <- function(data, withDsName=TRUE) {
  res <- lapply(names(data), function(name) {
    res <- split(data[[name]], as.factor(data[[name]]$id))
    if (withDsName) {
      names(res) <- paste(name, ": ", names(res))
    }
    res
  })
  names(res) <- names(data)
  res
}

splitBySeason <- function(data) {
  res <- lapply(names(data), function(name) {
    yd <- as.POSIXlt(data[[name]]$datetime)$yday
    season <- character(length(yd))
    season[yd < 74 | yd >= 319] = "Winter (15/11-14/3)"
    season[yd >= 74 & yd < 121] = "Spring (15/3-30/4)"
    season[yd >= 121 & yd < 274] = "Summer (1/5-30/9)"
    season[yd >= 274 & yd < 319] = "Autumn (1/10-14/11)"
    res <- split(data[[name]], as.factor(season))
    names(res) <- paste(name, ": ", names(res))
    res
  })
  names(res) <- names(data)
  res
}

splitByMonth <- function(data) {
  res <- lapply(names(data), function(name) {
    res <- split(data[[name]], as.factor(format(data[[name]]$datetime, "%m")))
    names(res) <- paste(name, ": ", names(res))
    res
  })
  names(res) <- names(data)
  res
}

plotDataQualityByStation <- function(data, dsName, colNames=c("lst","ground"), ids=c()) {
  sdata <- splitByStation(data, withDsName = FALSE)[[dsName]]
  if (length(ids) > 0) {
    sdata <- sdata[ids]
  }
  lapply(names(sdata), function(id) {
    iddata <- sdata[[id]]
    iddata <- iddata[order(iddata$datetime), ]
    newPage()
    par(mfcol=c(2,1), mar = c(2,4,2,1))
    plot(iddata$datetime, iddata[[colNames[1]]], type="line", col="#777777",
         main=sprintf("%s: %s (obs: %d)", dsName, id, nrow(iddata)),
         xlab=c(), ylab = "T & LST")
    points(iddata$datetime, iddata[[colNames[2]]], cex=0.5, col="red", pch=16)
    plot(iddata$datetime, iddata[[colNames[1]]] - iddata[[colNames[2]]], type="bar", pch=".", main=c(), xlab=c(), ylab="diff")
  })
}

densityPlot <- function(data, colNames=c("lst","ground"), fontFactor=1.0) {
  library(ggplot2)

  res <- lapply(names(data), function(name) {
    tt <- data[[name]][colNames]
    colnames(tt) <- c("lst", "ground")

    ft <- lm(ground ~ lst, data=tt)
    slope <- ft$coefficients[2]
    icept <- ft$coefficients[1]
    label <- sprintf("Obs: %d", nrow(tt))
    model <- sprintf("lm(%.2f+%.2f*LST,R^2=%.2f)", icept, slope, summary(ft)$r.squared)
    ggplot(tt, aes(x=ground,y=lst)) +
      xlim(-30, 30) + ylim(-30, 30) +
      labs(title=name, subtitle=label, caption=model, x="T", y="LST") +
      theme(
        plot.title = element_text(size = 11*fontFactor),
        plot.subtitle = element_text(size = 8*fontFactor),
        plot.caption = element_text(size = 7*fontFactor),
        axis.title = element_text(size = 9*fontFactor)
      ) +
      stat_density2d(aes(fill=..level..), geom="polygon", show.legend=FALSE) +
      scale_fill_gradient(low="#CCCCCC", high="blue") +
      geom_abline(slope=1, col="#777777", lwd=0.5) +
      geom_abline(slope=1/slope, intercept=-icept/slope, col="red")

  })
  names(res) <- names(data)
  res
}

plotTitlePage <- function(title, subtitle) {
  newPage()
  plot(0:10, type = "n", xaxt="n", yaxt="n", bty="n", xlab = "", ylab = "")
  text(5, 6, title, cex=2.0)
  text(5, 5, subtitle)
}

plot2x2Density <- function(data, colNames=c("lst","ground"), fontFactor=1.0) {
  p <- densityPlot(data, colNames=colNames, fontFactor=fontFactor)
  library(gridExtra)
  newPage()
  if (length(p) >= 4) {
    grid.arrange(p[[1]], p[[2]], p[[3]], p[[4]], nrow=2, ncol=2)
  } else {
    grid.arrange(p[[1]], p[[2]], p[[3]], nrow=2, ncol=2)
  }
}

plot3x4Density <- function(data, colNames=c("lst","ground"), fontFactor=1.0) {
  p <- densityPlot(data, colNames=colNames, fontFactor=fontFactor)
  library(gridExtra)
  newPage()
  grid.arrange(p[[12]], p[[1]], p[[2]], # FIXME: designed for months
               p[[3]], p[[4]], p[[5]],
               p[[6]], p[[7]], p[[8]],
               p[[9]], p[[10]], p[[11]], nrow=4, ncol=3)
}

plot4x5Density <- function(data, colNames=c("lst","ground"), fontFactor=1.0) {
  p <- densityPlot(data, colNames=colNames, fontFactor=fontFactor)
  library(gridExtra)
  newPage()
  grid.arrange(p[[1]], p[[2]], p[[3]], p[[4]],
               p[[5]], p[[6]], p[[7]], p[[8]],
               p[[9]], p[[10]], p[[11]], p[[12]],
               p[[13]], p[[14]], p[[15]], p[[16]],
               p[[17]], p[[18]], p[[19]], p[[20]], nrow=5, ncol=4)
}

to.time <- function(dt) {
  res <- (as.numeric(dt) - as.numeric(as.POSIXct(format(dt, "%Y-%m-%d"), format="%Y-%m-%d", tz="UTC")))/ 3600
  index <- res < 4
  res[index] = res[index] + 24
  res
}

plotTimeDistributions <- function(data) {
  newPage()
  par(mfcol=c(2,2))
  lapply(names(data), function(name) {
    tt <- data[[name]]
    t <- to.time(tt$datetime)
    hist(t, n=100,
      main=sprintf("Observation times: %s (UTC)", name),
      xlab="hour of day", ylab ="frequency")
  })
}
