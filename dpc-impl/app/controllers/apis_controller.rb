# frozen_string_literal: true

class ApisController < ApplicationController
  before_action :authenticate_user!

  def new
  end
end
