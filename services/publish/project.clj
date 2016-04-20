(use '[clojure.java.shell :only (sh)])
(require '[clojure.string :as string])

(defn git-ref
  []
  (or (System/getenv "GIT_COMMIT")
      (string/trim (:out (sh "git" "rev-parse" "HEAD")))
      ""))

(defproject org.iplantc/publish "5.2.6.0"
  :description "A pluggable REST API for submitting data publishing requests."
  :url "https://github.com/iPlantCollaborativeOpenSource/DE"
  :license {:name "BSD"
            :url "http://iplantcollaborative.org/sites/default/files/iPLANT-LICENSE.txt"}
  :manifest {"Git-Ref" ~(git-ref)}
  :uberjar-name "publish.jar"
  :dependencies [[org.clojure/clojure "1.7.0"]
                 [cheshire "5.5.0"]
                 [clj-http "2.0.0"]
                 [me.raynes/fs "1.4.6"]
                 [metosin/compojure-api "0.24.2"]
                 [org.iplantc/clojure-commons "5.2.6.0"]
                 [org.iplantc/common-cfg "5.2.6.0"]
                 [org.iplantc/common-cli "5.2.6.0"]
                 [org.iplantc/common-swagger-api "5.2.6.0"]
                 [org.iplantc/service-logging "5.2.6.0"]
                 [ring/ring-core "1.4.0"]
                 [ring/ring-jetty-adapter "1.4.0"]]
  :plugins [[lein-ring "0.9.6"]]
  :profiles {:dev {:resource-paths ["conf/test"]}}
  ;; compojure-api route macros should not be AOT compiled:
  ;; https://github.com/metosin/compojure-api/issues/135#issuecomment-121388539
  ;; https://github.com/metosin/compojure-api/issues/102
  :aot [#"publish.(?!routes).*"]
  :main publish.core
  :ring {:handler publish.routes/app
         :init    publish.core/init-service
         :port    31311}
  :uberjar-exclusions [#"(?i)META-INF/[^/]*[.](SF|DSA|RSA)"]
  :jvm-opts ["-Dlogback.configurationFile=/etc/iplant/de/logging/publish-logging.xml"])
