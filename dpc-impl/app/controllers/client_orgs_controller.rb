# frozen_string_literal: true

class ClientOrgsController < ApplicationController
  before_action :authenticate_user!
  
  def new
    @npi = generate_npi
  end

  def create
    binding.pry
    api_request = api_service.create_client_org()
  end

  private

  def api_service
    @api_service ||= ApiClient.new
  end
end
