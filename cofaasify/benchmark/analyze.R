library(tidyverse)
library(ggplot2)
library(stringi)
library(Cairo)
library(knitr)

#library(usethis)

#usethis::use_import_from("dplyr", "%>%")

do_stats <- function(filename, limit = 512) {
    data <- read_csv(filename) %>%
        group_by(type, size, repeats)

    res <- data %>%
        summarise(
            sd = sd(lat),
            `Latency (us)` = mean(lat)
        ) %>%
        filter(size <= limit) %>%
        mutate(
            size = as.factor(size),
            repeats = as.factor(repeats),
            type = stri_trans_totitle(type)
        ) %>%
        rename(
            `Size (KB)` = size,
            Repeats = repeats
        )
    res %>% ungroup()
}

read_data <- function(filename) {
    data <- read_csv(filename) %>%
        group_by(type, size, repeats)

    data %>%
        filter(size < 33) %>%
        mutate(
            size = as.factor(size),
            repeats = as.factor(repeats),
        )
}

calc_speedup <- function(t, a, b, cat, limit) {
    tt <- t %>%
        mutate(`Size int` = as.numeric(levels(`Size (KB)`))[`Size (KB)`]) %>%
        filter(`Size int` <= limit) %>%
        select(-`Size int`)
    type_a <- tt %>%
        filter(type == a) %>%
        rename(lat_a = `Latency (us)`) %>%
        mutate(type = cat) %>%
        select(-sd)

    type_b <- tt %>%
        filter(type == b) %>%
        rename(lat_b = `Latency (us)`)  %>%
        mutate(type = cat) %>%
        select(-sd)

    left_join(type_a, type_b) %>%
        mutate(Speedup = lat_b/lat_a) %>%
        select(-lat_a, -lat_b)
}




res <- do_stats("combined_baseline.csv")

do_stats("combined_write.csv")

res_go <- do_stats("combined_nogc2.csv")

res_rust <- do_stats("combined_rust.csv")


#res <- bind_rows(res_go, res_rust)

res <- do_stats("combined_otherbox.csv.xz")
rust <- calc_speedup(res, "Cofaas-Rust", "Native-Rust", "Rust", 512)
go <- calc_speedup(res, "Cofaas-Go-Nogc", "Native-Go-Nogc", "Go", 16)
go_gc <- calc_speedup(res, "Native-Go-Gccmp", "Native-Go-Nogc", "Go GC impact", 64)
go_rust <- calc_speedup(res, "Native-Rust", "Native-Go", "Go Rust impact native", 512)
wasm_go_rust <- calc_speedup(res, "Cofaas-Rust", "Cofaas-Go-Nogc","Go Rust Impact WASM", 16)
speedups <- bind_rows(rust, go)

res_latency <- do_stats("combined_latency_otherbox.csv.xz")
rust_latency <- calc_speedup(res_latency, "Latency-Cofaas-Rust", "Latency-Native-Rust", "Rust", 512)
go_latency <- calc_speedup(res_latency, "Latency-Cofaas-Go-Nogc", "Latency-Native-Go-Nogc", "Go", 16)
latency_speedups <- bind_rows(rust_latency, go_latency)

#res %>% ungroup() %>% select(type) %>% unique()

rawlat <- res %>%
    filter(Repeats == 20, (type == "Native-Rust" | type == "Cofaas-Rust")) %>%
    mutate(type = ifelse(type == "Native-Rust", "Native", "CoFaaS")) %>%
    select(-Repeats, -sd) %>%
    transmute(type, `Size (KB)`, latency = `Latency (us)`/1000)

print(speedups, n = 1000)

.gen_pdf_cairo <- function(plot, file_name) {
    file_name_pdf <- paste0(file_name, ".pdf")

    CairoPDF(file = file_name_pdf)
    print(plot)
    dev.off()
}

.gen_png_cairo <- function(plot, file_name) {
    file_name_pdf <- paste0(file_name, ".png")

    CairoPNG(file = file_name_pdf, dpi = 200)
    print(plot)
    dev.off()
}

#speedups

ggplot(res, aes(fill = Repeats, y = `Latency (us)`, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    geom_errorbar(aes(ymin = `Latency (us)` - sd, ymax = `Latency (us)` + sd),
                  width = .4,
                  linewidth = .2,
        position = position_dodge(.9
    )) +
    facet_wrap(~type) +
    labs(#title = "Roundtrip latency for requests issued to a
         #two-function DAG executed both natively and optimized using
         #the CoFaaS methodology for varying number of intra-DAG
                                        #requests repetitions",
        #title = "Request round-trip times (No GC)",
        x = "Intra-DAG request payload syze (KB)") +
       # theme_bw() +
    theme(aspect.ratio = 9/13,
          text =  element_text(size = 8),
          plot.margin=grid::unit(c(0,0,0,0), "mm")
          )

ggplot(rawlat, aes(fill = type, y = latency, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    labs(y = "Round-trip Latency (ms)") +
    theme(aspect.ratio = 5/13,
          text =  element_text(size = 45),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.10, .89),
          legend.title = element_blank())
ggsave("rawlat.pdf", device = cairo_pdf)
knitr::plot_crop("rawlat.pdf")


ggplot(speedups, aes(fill = Repeats, y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    annotate("rect", xmin=3.5, xmax=4.5, ymin=0 , ymax=6, alpha=0.2, color="blue", fill="blue") +
    #theme_gray(base_size = 50) +
    scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    ## geom_errorbar(aes(ymin = `Latency (us)` - sd, ymax = `Latency (us)` + sd),
    ##               width = .4,
    ##               linewidth = .2,
    ##     position = position_dodge(.9
    ## )) +
    facet_wrap(~type, scales = "free_x") +
    labs(#title = "Roundtrip latency for requests issued to a
         #two-function DAG executed both natively and optimized using
         #the CoFaaS methodology for varying number of intra-DAG
                                        #requests repetitions",
        #title = "Request round-trip times (No GC)",
        x = "Inter-function request payload syze (KB)") +
       # theme_bw() +
    theme(aspect.ratio = 9/13,
          text =  element_text(size = 15),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.95, .89))
ggsave("rust-go.pdf", device = cairo_pdf)
knitr::plot_crop("rust-go.pdf")

ggplot(speedups, aes(fill = Repeats, y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    annotate("rect", xmin=3.5, xmax=4.5, ymin=0 , ymax=6, alpha=0.2, color="blue", fill="blue") +
    scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    facet_wrap(~type, scales = "free_x") +
    labs(x = "Inter-Function request payload size (KB)",
         y = "Speedup Native to CoFaaS") +
    theme(aspect.ratio = 9/13,
          text =  element_text(size = 20),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.95, .89))
ggsave("rust-go.pdf", device = cairo_pdf)
knitr::plot_crop("rust-go.pdf")

ggplot(latency_speedups, aes(y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    annotate("rect", xmin=3.5, xmax=4.5, ymin=0 , ymax=100, alpha=0.2, color="blue", fill="blue") +
    #scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    facet_wrap(~type, scales = "free_x") +
    labs(x = "Inter-function request payload size (KB)", y = "Inter-function request latency speedup") +
    theme(aspect.ratio = 7/13,
          text =  element_text(size = 20),
          plot.margin=grid::unit(c(0,0,0,0), "mm"))
ggsave("rust-go-latency.pdf", device = cairo_pdf)
knitr::plot_crop("rust-go-latency.pdf")

ggplot(go_rust, aes(fill = Repeats, y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    labs(x = "Inter-Function request payload size (KB)",
         y = "Speedup from Go to Rust natively compiled") +
    theme(aspect.ratio = 7/13,
          text =  element_text(size = 45),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.10, .89))
ggsave("rust-go-native.pdf", device = cairo_pdf)
knitr::plot_crop("rust-go-native.pdf")

ggplot(wasm_go_rust, aes(fill = Repeats, y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    labs(x = "Inter-Function request payload size (KB)",
         y =  "Speedup from Go to Rust compiled to WASM") +
    theme(aspect.ratio = 7/13,
          text =  element_text(size = 45),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.10, .89))
ggsave("rust-go-wasm.pdf", device = cairo_pdf)
knitr::plot_crop("rust-go-wasm.pdf")

ggplot(go_gc, aes(fill = Repeats, y = Speedup, x = `Size (KB)`)) +
    geom_bar(position = "dodge", stat = "identity") +
    scale_y_continuous(breaks = c(1,2,3,4,5,6)) +
    labs(x = "Inter-Function request payload size (KB)",
         y = "Speedup (Go GC Disabled)") +
    theme(aspect.ratio = 4/13,
          text =  element_text(size = 45),
          plot.margin=grid::unit(c(0,0,0,0), "mm"),
          legend.position = c(.08, .25)) +
    guides(colour = guide_legend(override.aes = list(alpha = 1)))
ggsave("go-gc.pdf", device = cairo_pdf)
knitr::plot_crop("go-gc.pdf")


ggsave("plot.png")
#knitr::plot_crop()


## .gen_pdf_cairo(p, "plot")
.gen_png_cairo(p, "plot")
