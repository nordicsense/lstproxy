library(matching)

what <- "Terra-Day"

x <- read.matching(sprintf("/data/Research/Analysis/Fin-Meteo/%s.csv", what))
x <- x[is.finite(x$t) & is.finite(x$LST),]

st <- get.stations()
i <- match(x$id, st$id)

xx <- split(x, st$location[i])

matching.summary <- function(model, tt) {
  p <- predict(model, tt, interval = 'confidence', level = 0.99)
  byid <- sapply(test_ids, function(id) {
    i <- tt$id == id
    tryCatch(cor(tt$t[i], p[i,1], use="complete.obs")^2, error=function(...) NA)
  })
  byyear <- sapply(test_years, function(y) {
    i <- as.integer(format(tt$datetime, format = "%Y")) == y
    tryCatch(cor(tt$t[i], p[i,1], use="complete.obs")^2, error=function(...) NA)
  })
  list(ids=byid, years=byyear)
}

modelin <- fit.poly(xx$`in`, 2)
modeltown <- fit.poly(xx$town, 2)
modelcoastalland <-  fit.poly(xx[["coastal on land"]], 2)
modelcoastalisland <-  fit.poly(xx[["coastal on island"]], 2)

pdf(sprintf("/data/Research/Analysis/Fin-Meteo/%s-correlations.pdf", what), width=11, height=7)
par(mar=c(4,4,2,2))

plot.matching(modelcoastalisland, xx[["coastal on island"]], fontFactor = 2, newpage=FALSE, main="Heterogeneous: island")
plot.matching(modelcoastalland, xx[["coastal on land"]], fontFactor = 2, newpage=FALSE, main="Heterogeneous: coastal")
plot.matching(modeltown, xx$town, fontFactor = 2, newpage=FALSE, main="Heterogeneous: settlement")
plot.matching(modelin, xx$`in`, fontFactor = 2, newpage=FALSE, main="Homogeneous: inland")
dev.off()

# modelcoastalisland: 3.5893193    0.8917180   -0.0003174
# modelcoastalland: 3.654321     0.901595    -0.001476
# modelin: 3.754111     0.883413    -0.002325
# modeltown: 3.879162     0.868228    -0.002176

input <- data.frame(LST=seq(-25, 25, 0.5))
cols <- c("#7b3dff", "#bb4fd1", "#fa7d1e", "#67ed1a")
legend <- c("heterogenous: island", "heterogenous: coastal", "heterogenous: settlement", "homogeneous: inland")

pdf(sprintf("/data/Research/Analysis/Fin-Meteo/%s-by-location.pdf", what), width=7, height=11)
par(mfrow=c(2,1), mar=c(4,4,2,2))

plot(input$LST, predict(modelcoastalisland, input), type="l", col=cols[1], lwd=2, xlab="LST", ylab="Tp", ylim=c(-25,25))
lines(input$LST, predict(modelcoastalland, input), col=cols[2], lwd=2)
lines(input$LST, predict(modeltown, input), col=cols[3], lwd=2)
lines(input$LST, predict(modelin, input), col=cols[4], lwd=2)
abline(a=0,b=1, col="#CCCCCC")
legend("topleft", legend = legend, col=cols, lwd=2)

plot(input$LST, predict(modelcoastalisland, input)-input$LST, type="l", col=cols[1], lwd=2, xlab="LST", ylab="Tp - LST", ylim=c(-5,6))
lines(input$LST, predict(modelcoastalland, input)-input$LST, col=cols[2], lwd=2)
lines(input$LST, predict(modeltown, input)-input$LST, col=cols[3], lwd=2)
lines(input$LST, predict(modelin, input)-input$LST, col=cols[4], lwd=2)
legend("topright", legend = legend, col=cols, lwd=2)
abline(a=0,b=0, col="#CCCCCC")
dev.off()
