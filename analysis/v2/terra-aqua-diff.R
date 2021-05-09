

reshape <- function(data) {
  data$datetime <- substring(data$datetime, 1, 10)
  index <- !is.nan(data$LST) & !is.nan(data$t)
  data <- data[index, 1:4]
  split(data, data$id)
}

adata <- reshape(read.csv("/data/Research/Analysis/Fin-Meteo/Aqua-Day.csv"))
tdata <- reshape(read.csv("/data/Research/Analysis/Fin-Meteo/Terra-Day.csv"))

data <- lapply(adata, function(x) {
  id <- x$id[1]
  y <- tdata[[as.character(id)]]
  index <- match(x$datetime, y$datetime)
  x <- cbind(x, y[index, c(3, 4)])
  colnames(x) <- c("id", "datetime", "LST.a", "t.a", "LST.t", "t.t")
  x[!is.na(x$LST.t), ]
})

data <- do.call("rbind", data)
rownames(data) <- NULL

data <- split(data, substring(data$datetime, 6, 7))

res <- do.call("rbind", lapply(data, function(x) {
  c(mean(x$LST.a - x$LST.t), mean(x$t.a-x$t.t))
}))

colnames(res) <- c("LST(a-t)", "T(a-t)")

# Night
#LST(a-t)     T(a-t)
#01 -0.07449709  0.2650721
#02 -0.89544404 -0.9890607
#03 -2.05593302 -2.3548778
#04 -2.36906710 -3.0212222
#05 -3.34594945 -4.3639737
#06 -3.56433883 -4.5032429
#07 -3.24287043 -4.1931616
#08 -2.12086951 -2.7306492
#09 -1.13517506 -1.4457926
#10 -0.38900097 -0.4957191
#11  0.07200595  0.3018727
#12  0.02279790  0.7168877

#Day
#LST(a-t)     T(a-t)
#01  0.09190524 0.26186560
#02  0.62592224 0.77461120
#03  0.62144994 0.68080706
#04  0.44788058 0.46426886
#05  0.49964925 0.43197409
#06  0.58148956 0.37619439
#07  0.60906154 0.35983985
#08  0.52845000 0.37199769
#09  0.41534117 0.37889724
#10  0.35765424 0.36281560
#11  0.12178518 0.15441359
#12 -0.04858696 0.03794341
