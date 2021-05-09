library(matching)

x <- read.matching("/data/Research/Analysis/Fin-Meteo/Terra-Night.csv")

x <- x[is.finite(x$t) & is.finite(x$LST),]

test_years <- c(2019, 2018, 2014, 2010, 2006)
test_ids <- c(101398, 101009, 103794, 852678, 874863, 101689, 101831, 101246, 102005, 101237, 101339, 101028,
              101151, 101130, 101436, 101885, 101196, 101486, 101231, 101520, 100971, 101636, 101784)

i <- !(x$id %in% test_ids)
j <- !(as.integer(format(x$datetime, format = "%Y")) %in% test_years)

model <- fit.poly(x[i & j,], 2)

stats <- function(tt) {
  print(sort(unique(tt$id)))
  print(sort(unique(as.integer(format(tt$datetime, format = "%Y")))))
}


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
res <- matching.summary(model, x[!i & !j,])

plot.matching(model, x[!i & !j,], fontFactor = 2, newpage=TRUE)
stats(x[!i & !j,])

dev.new()
boxplot(res)

nrow(x)
nrow(x[i & j,])
nrow(x[!i & !j,])
