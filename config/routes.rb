Rails.application.routes.draw do
  resources :services do
    resources :machines
  end
  resources :machines
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
end
