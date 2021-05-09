fit.poly <- function(data, order=3) {
  # model <- lm(t ~ poly(LST, order), data=data)
  model <- lm(t ~ LST + I(LST^2), data=data)
  model
}
