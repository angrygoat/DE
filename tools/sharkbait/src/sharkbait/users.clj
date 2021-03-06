(ns sharkbait.users
  (:require [clojure.string :as string]
            [sharkbait.consts :as consts]
            [sharkbait.db :as db]
            [sharkbait.roles :as roles]
            [sharkbait.subjects :as subjects]))

(defn register-de-users
  [session de-users-folder subjects]
  (let [de-users-role (roles/create-role session de-users-folder consts/de-users-role-name)]
    (roles/replace-members de-users-role subjects)))
