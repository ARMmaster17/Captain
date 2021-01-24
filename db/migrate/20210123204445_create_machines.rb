class CreateMachines < ActiveRecord::Migration[6.1]
  def change
    create_table :machines do |t|
      t.string :hostname
      t.string :ip_address
      t.integer :vmid
      t.integer :cpu
      t.integer :ram
      t.integer :disk

      t.timestamps
    end
  end
end
