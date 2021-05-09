plot.matching <- function(model, data, fontFactor = 1.0, newpage = TRUE, main = "") {

  tt <- data.frame(T = data$t, LST = data$LST)
  tt <- tt[order(tt$LST),]
  tt <- tt[is.finite(tt$T) & is.finite(tt$LST),]

  p <- predict(model, tt, interval = 'confidence', level = 0.99)
  tt$Tp <- p[, 1]

  original <- ggplot(tt, aes(x = T, y = LST)) +
    xlim(-20, 30) +
    ylim(-20, 30) +
    labs(title = main, subtitle = "LST~T", caption = sprintf("R^2=%.2f", cor(tt$Tp, tt$T, use = "complete.obs")^2)) +
    theme(
      plot.title = element_text(size = 11 * fontFactor),
      plot.subtitle = element_text(size = 8 * fontFactor),
      plot.caption = element_text(size = 7 * fontFactor),
      axis.title = element_text(size = 9 * fontFactor)
    ) +
    stat_density2d(aes(fill = ..level..), geom = "polygon", show.legend = FALSE) +
    scale_fill_gradient(low = "#CCCCCC", high = "blue") +
    geom_abline(slope = 1, col = "#777777", lwd = 0.5) +
    geom_line(aes(x = Tp, y = LST), tt, col = "red")

  ft <- lm(T ~ Tp, data = tt)
  slope <- ft$coefficients[2]
  icept <- ft$coefficients[1]

  fitted <- ggplot(tt, aes(x = T, y = Tp)) +
    xlim(-20, 30) +
    ylim(-20, 30) +
    labs(title = main, subtitle = "Tp~T", caption = sprintf("R^2=%.2f", summary(ft)$r.squared)) +
    theme(
      plot.title = element_text(size = 11 * fontFactor),
      plot.subtitle = element_text(size = 8 * fontFactor),
      plot.caption = element_text(size = 7 * fontFactor),
      axis.title = element_text(size = 9 * fontFactor)
    ) +
    stat_density2d(aes(fill = ..level..), geom = "polygon", show.legend = FALSE) +
    scale_fill_gradient(low = "#CCCCCC", high = "blue") +
    geom_abline(slope = 1, col = "#777777", lwd = 0.5) +
    geom_abline(slope = 1 / slope, intercept = -icept / slope, col = "red")

  if (newpage) {
    dev.new()
  }
  grid.arrange(original, fitted, ncol = 2, nrow = 1)
}
