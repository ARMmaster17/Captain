class CreateServices < ActiveRecord::Migration[6.1]
  def change
    create_table :services do |t|
      t.string :name
      t.integer :scale
      t.integer :cpu
      t.integer :ram
      t.integer :disk
      t.string :hostname
      t.string :domain

      t.timestamps
    end

    add_reference :machines, :service
  end
end
