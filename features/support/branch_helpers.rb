# Returns the name of the branch that is currently checked out
def current_branch_name
  run("git branch").fetch(:out)
                        .split("\n")
                        .map(&:strip)
                        .select{|b| b[/^\*/]}
                        .first
                        .slice(2, 100)
end


# Returns the names of all existing local branches.
#
# Does not return the "master" branch nor remote branches.
#
# The branches are ordered this ways:
# * main branch
# * feature branches ordered alphabetically
def existing_local_branches
  actual_branches = run("git branch").fetch(:out)
                                          .split("\n")
                                          .map(&:strip)
                                          .map{|s| s.sub('* ', '')}
  actual_branches.delete('master')
  actual_main_branch = actual_branches.delete 'main'
  [actual_main_branch].concat(actual_branches)
                      .compact
end


def remote_branch_exists branch_name
  run("git branch -a | grep remotes/origin/#{branch_name} | wc -l")[:out] != '0'
end
