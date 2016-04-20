(ns publish.routes
  (:use [service-logging.middleware :only [wrap-logging add-user-to-context clean-context]]
        [clojure-commons.query-params :only [wrap-query-params]]
        [compojure.core :only [wrap-routes]]
        [common-swagger-api.schema]
        [ring.middleware.keyword-params :only [wrap-keyword-params]])
  (:require [compojure.route :as route]
            [cheshire.core :as json]
            [clojure-commons.exception :as cx]
            [publish.routes.status :as status-routes]
            [publish.util.config :as config]))

(defapi app
  {:exceptions cx/exception-handlers}
  (swagger-ui config/docs-uri
    :validator-url nil)
  (swagger-docs
   {:info {:title       (:desc config/svc-info)
           :description "Documentation for the publish API"
           :version     "2.0.0"}
    :tags [{:name "service-info", :description "Service Status Information"}]})
  (middlewares
   [clean-context
    wrap-keyword-params
    wrap-query-params
   (wrap-routes wrap-logging)]
   (context* "/" []
    :tags ["service-info"]
    status-routes/status))
  (middlewares
   [clean-context
    wrap-keyword-params
    wrap-query-params
    add-user-to-context
    wrap-logging]
   (route/not-found (json/encode {:success false :msg "unrecognized service path"}))))
