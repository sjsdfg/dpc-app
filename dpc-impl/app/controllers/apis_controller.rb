# frozen_string_literal: true

class ApisController < ApplicationController
  before_action :authenticate_user!

  def new
    @add_client_org
  end
end
